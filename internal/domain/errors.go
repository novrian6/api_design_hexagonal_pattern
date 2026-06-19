// internal/domain/errors.go
package domain

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrEmployeeNotFound   = errors.New("employee not found")
	ErrEmptyID            = errors.New("employee id cannot be empty")
	ErrEmptyName          = errors.New("employee name cannot be empty")
	ErrNegativeSalary     = errors.New("salary cannot be negative")
	ErrEmptyEmail         = errors.New("employee email cannot be empty")
)
