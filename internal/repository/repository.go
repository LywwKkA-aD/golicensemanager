package repository

import (
	"context"

	"github.com/LywwKkA-aD/golicensemanager/internal/models"
	"github.com/google/uuid"
)

// ApplicationRepository handles database operations for applications
type ApplicationRepository interface {
	Create(ctx context.Context, app *models.Application) (*models.Application, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error)
	GetByAPIKey(ctx context.Context, apiKey string) (*models.Application, error)
	List(ctx context.Context) ([]models.Application, error)
	Update(ctx context.Context, app *models.Application) (*models.Application, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// LicenseTypeRepository handles database operations for license types
type LicenseTypeRepository interface {
	Create(ctx context.Context, licenseType *models.LicenseType) (*models.LicenseType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.LicenseType, error)
	List(ctx context.Context, applicationID uuid.UUID) ([]models.LicenseType, error)
	Update(ctx context.Context, licenseType *models.LicenseType) (*models.LicenseType, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// LicenseRepository handles database operations for licenses
type LicenseRepository interface {
	Create(ctx context.Context, license *models.License) (*models.License, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.License, error)
	GetByKey(ctx context.Context, licenseKey string) (*models.License, error)
	List(ctx context.Context, filters LicenseFilters) ([]models.License, error)
	Update(ctx context.Context, license *models.License) (*models.License, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CreateActivity(ctx context.Context, activity *models.LicenseActivity) error
	GetActivities(ctx context.Context, licenseID uuid.UUID) ([]models.LicenseActivity, error)
	HasActiveClientLicenses(ctx context.Context, applicationID, clientID uuid.UUID) (bool, error)
}

// ClientRepository handles database operations for clients
type ClientRepository interface {
	Create(ctx context.Context, client *models.Client) (*models.Client, error)
	GetByID(ctx context.Context, applicationID, id uuid.UUID) (*models.Client, error)
	List(ctx context.Context, filters ClientFilters) ([]models.Client, error)
	Update(ctx context.Context, client *models.Client) (*models.Client, error)
	Delete(ctx context.Context, applicationID, id uuid.UUID) error
	ExistsByEmail(ctx context.Context, applicationID uuid.UUID, email string) (bool, error)
}

// LicenseFilters defines the available filters for listing licenses
type LicenseFilters struct {
	ApplicationID uuid.UUID
	ClientID      *uuid.UUID
	IsActive      *bool
	IsRevoked     *bool
}

// ClientFilters defines the available filters for listing clients
type ClientFilters struct {
	ApplicationID uuid.UUID
	IsActive      *bool
	Search        string
}
