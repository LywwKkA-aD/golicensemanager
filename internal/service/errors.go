package service

import "errors"

var (
	// Common errors
	ErrNotFound     = errors.New("record not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")

	// Application specific errors
	ErrInvalidAPICredentials = errors.New("invalid API credentials")

	// License specific errors
	ErrLicenseInvalid            = errors.New("license is invalid")
	ErrLicenseExpired            = errors.New("license has expired")
	ErrLicenseRevoked            = errors.New("license has been revoked")
	ErrLicenseUsageLimitExceeded = errors.New("license usage limit exceeded")

	// Client specific errors
	ErrDuplicateEmail          = errors.New("email already exists")
	ErrClientHasActiveLicenses = errors.New("client has active licenses")

	// General business errors
	ErrInvalidDateRange = errors.New("invalid date range")
	ErrFutureDate       = errors.New("date cannot be in the future")
	ErrInvalidStatus    = errors.New("invalid status")
)
