package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/LywwKkA-aD/golicensemanager/internal/models"
	"github.com/LywwKkA-aD/golicensemanager/internal/repository"
)

type LicenseFilters struct {
	ApplicationID uuid.UUID  `form:"application_id"`
	ClientID      *uuid.UUID `form:"client_id"`
	IsActive      *bool      `form:"is_active"`
	IsRevoked     *bool      `form:"is_revoked"`
}

type ValidationResult struct {
	Valid     bool                   `json:"valid"`
	Message   string                 `json:"message"`
	ExpiresAt time.Time              `json:"expires_at"`
	Features  map[string]interface{} `json:"features"`
}

type LicenseService interface {
	Create(ctx context.Context, license *models.License) (*models.License, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.License, error)
	List(ctx context.Context, filters LicenseFilters) ([]models.License, error)
	Update(ctx context.Context, license *models.License) (*models.License, error)
	Revoke(ctx context.Context, id uuid.UUID, reason string) error
	Validate(ctx context.Context, licenseKey string) (*ValidationResult, error)
	GetByKey(ctx context.Context, licenseKey string) (*models.License, error)
	RecordActivity(ctx context.Context, activity *models.LicenseActivity) error
	CheckUsage(ctx context.Context, licenseKey string, usage map[string]interface{}) error
}

type licenseService struct {
	repo            repository.LicenseRepository
	licenseTypeRepo repository.LicenseTypeRepository
	logger          *zap.SugaredLogger
}

func NewLicenseService(
	repo repository.LicenseRepository,
	licenseTypeRepo repository.LicenseTypeRepository,
	logger *zap.SugaredLogger,
) LicenseService {
	return &licenseService{
		repo:            repo,
		licenseTypeRepo: licenseTypeRepo,
		logger:          logger,
	}
}

func (s *licenseService) Create(ctx context.Context, license *models.License) (*models.License, error) {
	// Validate input
	if err := validateLicense(license); err != nil {
		return nil, err
	}

	// Get license type to set duration
	licenseType, err := s.licenseTypeRepo.GetByID(ctx, license.LicenseTypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get license type: %w", err)
	}

	// Set license dates
	license.StartDate = time.Now()
	license.ExpiryDate = license.StartDate.AddDate(0, 0, licenseType.DurationDays)

	// Generate license key
	licenseKey, err := generateLicenseKey(license)
	if err != nil {
		return nil, fmt.Errorf("failed to generate license key: %w", err)
	}
	license.LicenseKey = licenseKey

	// Set initial usage limits from license type if not provided
	if license.UsageLimits == nil {
		license.UsageLimits = licenseType.Features
	}

	// Initialize current usage
	license.CurrentUsage = make(map[string]interface{})

	return s.repo.Create(ctx, license)
}

func (s *licenseService) GetByID(ctx context.Context, id uuid.UUID) (*models.License, error) {
	license, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return license, nil
}

func (s *licenseService) List(ctx context.Context, filters LicenseFilters) ([]models.License, error) {
	return s.repo.List(ctx, repository.LicenseFilters{
		ApplicationID: filters.ApplicationID,
		ClientID:      filters.ClientID,
		IsActive:      filters.IsActive,
		IsRevoked:     filters.IsRevoked,
	})
}

func (s *licenseService) Update(ctx context.Context, license *models.License) (*models.License, error) {
	// Validate input
	if err := validateLicense(license); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByID(ctx, license.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Preserve certain fields
	license.LicenseKey = existing.LicenseKey
	license.StartDate = existing.StartDate
	license.ExpiryDate = existing.ExpiryDate
	license.CurrentUsage = existing.CurrentUsage

	return s.repo.Update(ctx, license)
}

func (s *licenseService) Revoke(ctx context.Context, id uuid.UUID, reason string) error {
	license, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return ErrNotFound
		}
		return err
	}

	license.IsRevoked = true
	license.RevocationReason = &reason
	license.IsActive = false

	_, err = s.repo.Update(ctx, license)
	if err != nil {
		return err
	}

	// Record revocation activity
	activity := &models.LicenseActivity{
		LicenseID:    id,
		ActivityType: "revocation",
		Description:  fmt.Sprintf("License revoked: %s", reason),
		Metadata: map[string]interface{}{
			"reason": reason,
		},
	}

	return s.RecordActivity(ctx, activity)
}

func (s *licenseService) Validate(ctx context.Context, licenseKey string) (*ValidationResult, error) {
	license, err := s.repo.GetByKey(ctx, licenseKey)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrLicenseInvalid
		}
		return nil, err
	}

	result := &ValidationResult{
		Valid:     true,
		ExpiresAt: license.ExpiryDate,
		Features:  make(map[string]interface{}),
	}

	// Check if license is revoked
	if license.IsRevoked {
		result.Valid = false
		result.Message = "License has been revoked"
		return result, ErrLicenseRevoked
	}

	// Check if license is expired
	if time.Now().After(license.ExpiryDate) {
		result.Valid = false
		result.Message = "License has expired"
		return result, ErrLicenseExpired
	}

	// Check if license is active
	if !license.IsActive {
		result.Valid = false
		result.Message = "License is not active"
		return result, ErrLicenseInvalid
	}

	// Get license type to include features
	licenseType, err := s.licenseTypeRepo.GetByID(ctx, license.LicenseTypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get license type: %w", err)
	}
	result.Features = licenseType.Features

	// Update last check time
	now := time.Now()
	license.LastCheck = &now
	if _, err := s.repo.Update(ctx, license); err != nil {
		s.logger.Warnf("Failed to update license last check time: %v", err)
	}

	// Record validation activity
	activity := &models.LicenseActivity{
		LicenseID:    license.ID,
		ActivityType: "validation",
		Description:  "License validated successfully",
	}
	if err := s.RecordActivity(ctx, activity); err != nil {
		s.logger.Warnf("Failed to record license activity: %v", err)
	}

	return result, nil
}

func (s *licenseService) GetByKey(ctx context.Context, licenseKey string) (*models.License, error) {
	license, err := s.repo.GetByKey(ctx, licenseKey)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return license, nil
}

func (s *licenseService) RecordActivity(ctx context.Context, activity *models.LicenseActivity) error {
	if activity.CreatedAt.IsZero() {
		activity.CreatedAt = time.Now()
	}
	return s.repo.CreateActivity(ctx, activity)
}

func (s *licenseService) CheckUsage(ctx context.Context, licenseKey string, usage map[string]interface{}) error {
	license, err := s.GetByKey(ctx, licenseKey)
	if err != nil {
		return err
	}

	// Validate the license first
	if _, err := s.Validate(ctx, licenseKey); err != nil {
		return err
	}

	// Check each usage metric against limits
	for metric, value := range usage {
		limit, exists := license.UsageLimits[metric]
		if !exists {
			continue // Skip metrics that don't have limits
		}

		// Convert values to float64 for comparison
		currentValue, ok := value.(float64)
		if !ok {
			return fmt.Errorf("%w: invalid usage value type for metric %s", ErrInvalidInput, metric)
		}

		limitValue, ok := limit.(float64)
		if !ok {
			return fmt.Errorf("%w: invalid limit value type for metric %s", ErrInvalidInput, metric)
		}

		if currentValue > limitValue {
			return fmt.Errorf("%w: %s exceeds allowed limit", ErrLicenseUsageLimitExceeded, metric)
		}

		// Update current usage
		license.CurrentUsage[metric] = currentValue
	}

	// Update license with new usage data
	_, err = s.repo.Update(ctx, license)
	return err
}

// Helper functions

func validateLicense(license *models.License) error {
	if license.ApplicationID == uuid.Nil {
		return fmt.Errorf("%w: application ID is required", ErrInvalidInput)
	}
	if license.LicenseTypeID == uuid.Nil {
		return fmt.Errorf("%w: license type ID is required", ErrInvalidInput)
	}
	if license.ClientID == uuid.Nil {
		return fmt.Errorf("%w: client ID is required", ErrInvalidInput)
	}
	return nil
}

func generateLicenseKey(license *models.License) (string, error) {
	// Create a unique string combining multiple fields
	unique := fmt.Sprintf("%s-%s-%s-%d",
		license.ApplicationID.String(),
		license.ClientID.String(),
		license.LicenseTypeID.String(),
		time.Now().UnixNano(),
	)

	// Create SHA-256 hash
	hash := sha256.New()
	hash.Write([]byte(unique))

	// Return first 32 characters of hex-encoded hash
	return hex.EncodeToString(hash.Sum(nil))[:32], nil
}
