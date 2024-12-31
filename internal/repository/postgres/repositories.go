package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/LywwKkA-aD/golicensemanager/internal/models"
	"github.com/LywwKkA-aD/golicensemanager/internal/repository"
)

// applicationRepo implements repository.ApplicationRepository
type applicationRepo struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) repository.ApplicationRepository {
	return &applicationRepo{db: db}
}

func (r *applicationRepo) Create(ctx context.Context, app *models.Application) (*models.Application, error) {
	if err := r.db.WithContext(ctx).Create(app).Error; err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}
	return app, nil
}

func (r *applicationRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	var app models.Application
	if err := r.db.WithContext(ctx).First(&app, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	return &app, nil
}

func (r *applicationRepo) GetByAPIKey(ctx context.Context, apiKey string) (*models.Application, error) {
	var app models.Application
	if err := r.db.WithContext(ctx).First(&app, "api_key = ?", apiKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get application by API key: %w", err)
	}
	return &app, nil
}

func (r *applicationRepo) List(ctx context.Context) ([]models.Application, error) {
	var apps []models.Application
	if err := r.db.WithContext(ctx).Find(&apps).Error; err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}
	return apps, nil
}

func (r *applicationRepo) Update(ctx context.Context, app *models.Application) (*models.Application, error) {
	if err := r.db.WithContext(ctx).Save(app).Error; err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}
	return app, nil
}

func (r *applicationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Application{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete application: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// licenseTypeRepo implements repository.LicenseTypeRepository
type licenseTypeRepo struct {
	db *gorm.DB
}

func NewLicenseTypeRepository(db *gorm.DB) repository.LicenseTypeRepository {
	return &licenseTypeRepo{db: db}
}

func (r *licenseTypeRepo) Create(ctx context.Context, licenseType *models.LicenseType) (*models.LicenseType, error) {
	if err := r.db.WithContext(ctx).Create(licenseType).Error; err != nil {
		return nil, fmt.Errorf("failed to create license type: %w", err)
	}
	return licenseType, nil
}

func (r *licenseTypeRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.LicenseType, error) {
	var licenseType models.LicenseType
	if err := r.db.WithContext(ctx).First(&licenseType, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get license type: %w", err)
	}
	return &licenseType, nil
}

func (r *licenseTypeRepo) List(ctx context.Context, applicationID uuid.UUID) ([]models.LicenseType, error) {
	var types []models.LicenseType
	if err := r.db.WithContext(ctx).Where("application_id = ?", applicationID).Find(&types).Error; err != nil {
		return nil, fmt.Errorf("failed to list license types: %w", err)
	}
	return types, nil
}

func (r *licenseTypeRepo) Update(ctx context.Context, licenseType *models.LicenseType) (*models.LicenseType, error) {
	if err := r.db.WithContext(ctx).Save(licenseType).Error; err != nil {
		return nil, fmt.Errorf("failed to update license type: %w", err)
	}
	return licenseType, nil
}

func (r *licenseTypeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.LicenseType{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete license type: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// licenseRepo implements repository.LicenseRepository
type licenseRepo struct {
	db *gorm.DB
}

func NewLicenseRepository(db *gorm.DB) repository.LicenseRepository {
	return &licenseRepo{db: db}
}

func (r *licenseRepo) Create(ctx context.Context, license *models.License) (*models.License, error) {
	if err := r.db.WithContext(ctx).Create(license).Error; err != nil {
		return nil, fmt.Errorf("failed to create license: %w", err)
	}
	return license, nil
}

func (r *licenseRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.License, error) {
	var license models.License
	if err := r.db.WithContext(ctx).
		Preload("LicenseType").
		Preload("Client").
		First(&license, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get license: %w", err)
	}
	return &license, nil
}

func (r *licenseRepo) GetByKey(ctx context.Context, licenseKey string) (*models.License, error) {
	var license models.License
	if err := r.db.WithContext(ctx).
		Preload("LicenseType").
		Preload("Client").
		First(&license, "license_key = ?", licenseKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get license by key: %w", err)
	}
	return &license, nil
}

func (r *licenseRepo) List(ctx context.Context, filters repository.LicenseFilters) ([]models.License, error) {
	var licenses []models.License
	query := r.db.WithContext(ctx).
		Preload("LicenseType").
		Preload("Client").
		Where("application_id = ?", filters.ApplicationID)

	if filters.ClientID != nil {
		query = query.Where("client_id = ?", *filters.ClientID)
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	if filters.IsRevoked != nil {
		query = query.Where("is_revoked = ?", *filters.IsRevoked)
	}

	if err := query.Find(&licenses).Error; err != nil {
		return nil, fmt.Errorf("failed to list licenses: %w", err)
	}
	return licenses, nil
}

func (r *licenseRepo) Update(ctx context.Context, license *models.License) (*models.License, error) {
	if err := r.db.WithContext(ctx).Save(license).Error; err != nil {
		return nil, fmt.Errorf("failed to update license: %w", err)
	}
	return license, nil
}

func (r *licenseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.License{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete license: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *licenseRepo) CreateActivity(ctx context.Context, activity *models.LicenseActivity) error {
	if err := r.db.WithContext(ctx).Create(activity).Error; err != nil {
		return fmt.Errorf("failed to create license activity: %w", err)
	}
	return nil
}

func (r *licenseRepo) GetActivities(ctx context.Context, licenseID uuid.UUID) ([]models.LicenseActivity, error) {
	var activities []models.LicenseActivity
	if err := r.db.WithContext(ctx).
		Where("license_id = ?", licenseID).
		Order("created_at DESC").
		Find(&activities).Error; err != nil {
		return nil, fmt.Errorf("failed to get license activities: %w", err)
	}
	return activities, nil
}

func (r *licenseRepo) HasActiveClientLicenses(ctx context.Context, applicationID, clientID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.License{}).
		Where("application_id = ? AND client_id = ? AND is_active = ? AND is_revoked = ?",
			applicationID, clientID, true, false).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check active licenses: %w", err)
	}
	return count > 0, nil
}

// clientRepo implements repository.ClientRepository
type clientRepo struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) repository.ClientRepository {
	return &clientRepo{db: db}
}

func (r *clientRepo) Create(ctx context.Context, client *models.Client) (*models.Client, error) {
	if err := r.db.WithContext(ctx).Create(client).Error; err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return client, nil
}

func (r *clientRepo) GetByID(ctx context.Context, applicationID, id uuid.UUID) (*models.Client, error) {
	var client models.Client
	if err := r.db.WithContext(ctx).
		Where("application_id = ? AND id = ?", applicationID, id).
		First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}
	return &client, nil
}

func (r *clientRepo) List(ctx context.Context, filters repository.ClientFilters) ([]models.Client, error) {
	var clients []models.Client
	query := r.db.WithContext(ctx).Where("application_id = ?", filters.ApplicationID)

	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		query = query.Where(
			"name ILIKE ? OR email ILIKE ? OR company ILIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if err := query.Find(&clients).Error; err != nil {
		return nil, fmt.Errorf("failed to list clients: %w", err)
	}
	return clients, nil
}

func (r *clientRepo) Update(ctx context.Context, client *models.Client) (*models.Client, error) {
	if err := r.db.WithContext(ctx).Save(client).Error; err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}
	return client, nil
}

func (r *clientRepo) Delete(ctx context.Context, applicationID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("application_id = ? AND id = ?", applicationID, id).
		Delete(&models.Client{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete client: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *clientRepo) ExistsByEmail(ctx context.Context, applicationID uuid.UUID, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Client{}).
		Where("application_id = ? AND email = ?", applicationID, email).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}
