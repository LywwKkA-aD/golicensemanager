package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/LywwKkA-aD/golicensemanager/internal/models"
	"github.com/LywwKkA-aD/golicensemanager/internal/repository"
)

type ApplicationService interface {
	Create(ctx context.Context, app *models.Application) (*models.Application, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error)
	List(ctx context.Context) ([]models.Application, error)
	Update(ctx context.Context, app *models.Application) (*models.Application, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GenerateToken(ctx context.Context, apiKey, apiSecret string) (string, error)
	ValidateAPICredentials(ctx context.Context, apiKey, apiSecret string) (*models.Application, error)
}

type applicationService struct {
	repo   repository.ApplicationRepository
	logger *zap.SugaredLogger
}

func NewApplicationService(repo repository.ApplicationRepository, logger *zap.SugaredLogger) ApplicationService {
	return &applicationService{
		repo:   repo,
		logger: logger,
	}
}

func (s *applicationService) Create(ctx context.Context, app *models.Application) (*models.Application, error) {
	// Generate API credentials
	apiKey, err := generateSecureKey(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	apiSecret, err := generateSecureKey(64)
	if err != nil {
		return nil, fmt.Errorf("failed to generate API secret: %w", err)
	}

	app.APIKey = apiKey
	app.APISecret = apiSecret

	// Validate input
	if err := validateApplication(app); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, app)
}

func (s *applicationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	app, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return app, nil
}

func (s *applicationService) List(ctx context.Context) ([]models.Application, error) {
	return s.repo.List(ctx)
}

func (s *applicationService) Update(ctx context.Context, app *models.Application) (*models.Application, error) {
	// Validate input
	if err := validateApplication(app); err != nil {
		return nil, err
	}

	// Check existence
	existing, err := s.repo.GetByID(ctx, app.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Preserve API credentials
	app.APIKey = existing.APIKey
	app.APISecret = existing.APISecret

	return s.repo.Update(ctx, app)
}

func (s *applicationService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *applicationService) GenerateToken(ctx context.Context, apiKey, apiSecret string) (string, error) {
	app, err := s.ValidateAPICredentials(ctx, apiKey, apiSecret)
	if err != nil {
		return "", err
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"application_id": app.ID.String(),
		"api_key":        app.APIKey,
		"exp":            time.Now().Add(24 * time.Hour).Unix(),
	})

	// TODO: Get JWT secret from config
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

func (s *applicationService) ValidateAPICredentials(ctx context.Context, apiKey, apiSecret string) (*models.Application, error) {
	app, err := s.repo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrInvalidAPICredentials
		}
		return nil, err
	}

	if app.APISecret != apiSecret {
		return nil, ErrInvalidAPICredentials
	}

	return app, nil
}

// Helper functions

func generateSecureKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func validateApplication(app *models.Application) error {
	if app.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	return nil
}
