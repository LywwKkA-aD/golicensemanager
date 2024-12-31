package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/LywwKkA-aD/golicensemanager/internal/models"
	"github.com/LywwKkA-aD/golicensemanager/internal/repository"
)

type ClientFilters struct {
	ApplicationID uuid.UUID `form:"application_id"`
	IsActive      *bool     `form:"is_active"`
	Search        string    `form:"search"`
}

type ClientService interface {
	Create(ctx context.Context, client *models.Client) (*models.Client, error)
	GetByID(ctx context.Context, applicationID, id uuid.UUID) (*models.Client, error)
	List(ctx context.Context, filters ClientFilters) ([]models.Client, error)
	Update(ctx context.Context, client *models.Client) (*models.Client, error)
	Delete(ctx context.Context, applicationID, id uuid.UUID) error
	GetClientLicenses(ctx context.Context, applicationID, clientID uuid.UUID) ([]models.License, error)
}

type clientService struct {
	repo        repository.ClientRepository
	licenseRepo repository.LicenseRepository
	logger      *zap.SugaredLogger
}

func NewClientService(
	repo repository.ClientRepository,
	licenseRepo repository.LicenseRepository,
	logger *zap.SugaredLogger,
) ClientService {
	return &clientService{
		repo:        repo,
		licenseRepo: licenseRepo,
		logger:      logger,
	}
}

func (s *clientService) Create(ctx context.Context, client *models.Client) (*models.Client, error) {
	// Validate input
	if err := validateClient(client); err != nil {
		return nil, err
	}

	// Check for duplicate email within the same application
	exists, err := s.repo.ExistsByEmail(ctx, client.ApplicationID, client.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateEmail
	}

	// Set default values
	if client.Metadata == nil {
		client.Metadata = make(map[string]interface{})
	}
	client.IsActive = true

	return s.repo.Create(ctx, client)
}

func (s *clientService) GetByID(ctx context.Context, applicationID, id uuid.UUID) (*models.Client, error) {
	client, err := s.repo.GetByID(ctx, applicationID, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return client, nil
}

func (s *clientService) List(ctx context.Context, filters ClientFilters) ([]models.Client, error) {
	return s.repo.List(ctx, repository.ClientFilters{
		ApplicationID: filters.ApplicationID,
		IsActive:      filters.IsActive,
		Search:        filters.Search,
	})
}

func (s *clientService) Update(ctx context.Context, client *models.Client) (*models.Client, error) {
	// Validate input
	if err := validateClient(client); err != nil {
		return nil, err
	}

	// Check existence
	existing, err := s.repo.GetByID(ctx, client.ApplicationID, client.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Check for duplicate email if email is being changed
	if existing.Email != client.Email {
		exists, err := s.repo.ExistsByEmail(ctx, client.ApplicationID, client.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrDuplicateEmail
		}
	}

	// Preserve certain fields
	client.CreatedAt = existing.CreatedAt
	if client.Metadata == nil {
		client.Metadata = existing.Metadata
	}

	return s.repo.Update(ctx, client)
}

func (s *clientService) Delete(ctx context.Context, applicationID, id uuid.UUID) error {
	// Check if client has any active licenses
	hasActiveLicenses, err := s.licenseRepo.HasActiveClientLicenses(ctx, applicationID, id)
	if err != nil {
		return err
	}
	if hasActiveLicenses {
		return ErrClientHasActiveLicenses
	}

	// Soft delete the client by setting IsActive to false
	client, err := s.repo.GetByID(ctx, applicationID, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return ErrNotFound
		}
		return err
	}

	client.IsActive = false
	_, err = s.repo.Update(ctx, client)
	return err
}

func (s *clientService) GetClientLicenses(ctx context.Context, applicationID, clientID uuid.UUID) ([]models.License, error) {
	// First verify the client exists
	if _, err := s.GetByID(ctx, applicationID, clientID); err != nil {
		return nil, err
	}

	// Get client's licenses
	return s.licenseRepo.List(ctx, repository.LicenseFilters{
		ApplicationID: applicationID,
		ClientID:      &clientID,
	})
}

// Helper functions

func validateClient(client *models.Client) error {
	if client.ApplicationID == uuid.Nil {
		return fmt.Errorf("%w: application ID is required", ErrInvalidInput)
	}
	if client.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if client.Email == "" {
		return fmt.Errorf("%w: email is required", ErrInvalidInput)
	}
	// Could add email format validation here
	return nil
}
