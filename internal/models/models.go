package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	CreatedAt time.Time `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Application struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Version     string    `gorm:"type:varchar(50)" json:"version"`
	APIKey      string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"api_key"`
	APISecret   string    `gorm:"type:varchar(128);not null" json:"api_secret"`
	Base
}

type LicenseType struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ApplicationID uuid.UUID      `gorm:"type:uuid;not null" json:"application_id"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	DurationDays  int            `gorm:"not null" json:"duration_days"`
	Price         float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	Features      map[string]any `gorm:"type:jsonb;default:'{}'" json:"features"`
	Application   Application    `gorm:"foreignKey:ApplicationID;constraint:OnDelete:CASCADE" json:"-"`
	Base
}

type Client struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ApplicationID uuid.UUID      `gorm:"type:uuid;not null" json:"application_id"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	Email         string         `gorm:"type:varchar(255);not null" json:"email"`
	Company       string         `gorm:"type:varchar(255)" json:"company"`
	ContactPerson string         `gorm:"type:varchar(255)" json:"contact_person"`
	Phone         string         `gorm:"type:varchar(50)" json:"phone"`
	Metadata      map[string]any `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	Application   Application    `gorm:"foreignKey:ApplicationID;constraint:OnDelete:CASCADE" json:"-"`
	Base
}

type License struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ApplicationID    uuid.UUID      `gorm:"type:uuid;not null" json:"application_id"`
	LicenseTypeID    uuid.UUID      `gorm:"type:uuid;not null" json:"license_type_id"`
	ClientID         uuid.UUID      `gorm:"type:uuid;not null" json:"client_id"`
	LicenseKey       string         `gorm:"type:varchar(128);uniqueIndex;not null" json:"license_key"`
	StartDate        time.Time      `gorm:"type:date;not null" json:"start_date"`
	ExpiryDate       time.Time      `gorm:"type:date;not null" json:"expiry_date"`
	UsageLimits      map[string]any `gorm:"type:jsonb;default:'{}'" json:"usage_limits"`
	CurrentUsage     map[string]any `gorm:"type:jsonb;default:'{}'" json:"current_usage"`
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	IsRevoked        bool           `gorm:"default:false" json:"is_revoked"`
	RevocationReason *string        `gorm:"type:text" json:"revocation_reason"`
	LastCheck        *time.Time     `gorm:"type:timestamp with time zone" json:"last_check"`
	Application      Application    `gorm:"foreignKey:ApplicationID;constraint:OnDelete:CASCADE" json:"-"`
	LicenseType      LicenseType    `gorm:"foreignKey:LicenseTypeID" json:"-"`
	Client           Client         `gorm:"foreignKey:ClientID" json:"-"`
	Base
}

type LicenseActivity struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	LicenseID    uuid.UUID      `gorm:"type:uuid;not null" json:"license_id"`
	ActivityType string         `gorm:"type:varchar(50);not null" json:"activity_type"`
	Description  string         `gorm:"type:text" json:"description"`
	Metadata     map[string]any `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	IPAddress    string         `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent    string         `gorm:"type:text" json:"user_agent"`
	License      License        `gorm:"foreignKey:LicenseID;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt    time.Time      `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type APIToken struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ApplicationID uuid.UUID      `gorm:"type:uuid;not null" json:"application_id"`
	Token         string         `gorm:"type:varchar(128);uniqueIndex;not null" json:"token"`
	Description   string         `gorm:"type:text" json:"description"`
	Permissions   map[string]any `gorm:"type:jsonb;default:'{}'" json:"permissions"`
	ExpiresAt     *time.Time     `gorm:"type:timestamp with time zone" json:"expires_at"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	Application   Application    `gorm:"foreignKey:ApplicationID;constraint:OnDelete:CASCADE" json:"-"`
	Base
}

// BeforeCreate will set a UUID rather than numeric ID
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	base.CreatedAt = time.Now()
	base.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate will update the updated_at timestamp
func (base *Base) BeforeUpdate(tx *gorm.DB) error {
	base.UpdatedAt = time.Now()
	return nil
}
