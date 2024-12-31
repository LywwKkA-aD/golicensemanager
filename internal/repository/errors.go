package repository

import "errors"

var (
	// ErrNotFound is returned when a requested record is not found
	ErrNotFound = errors.New("record not found")

	// ErrDuplicateKey is returned when trying to insert a record with a duplicate unique key
	ErrDuplicateKey = errors.New("duplicate key value")

	// ErrInvalidInput is returned when the input data is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrForeignKeyViolation is returned when a foreign key constraint is violated
	ErrForeignKeyViolation = errors.New("foreign key violation")
)
