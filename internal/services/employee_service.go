// internal/services/employee_service.go
package services

import (
	"context"
	"errors"
	"fmt"

	"decoupled_contract_services/internal/domain"
	"decoupled_contract_services/internal/ports"
)

// employeeServiceImpl adalah implementasi konkret dari EmployeeService (inbound port)
type employeeServiceImpl struct {
	repo ports.EmployeeRepository // Bergantung pada interface (outbound port)
}

// NewEmployeeService adalah constructor dengan dependency injection
func NewEmployeeService(repo ports.EmployeeRepository) ports.EmployeeService {
	return &employeeServiceImpl{repo: repo}
}

func (s *employeeServiceImpl) GetAllEmployees(ctx context.Context) ([]domain.Employee, error) {
	employees, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get all employees: %w", err)
	}
	return employees, nil
}

func (s *employeeServiceImpl) GetEmployeeByID(ctx context.Context, id string) (domain.Employee, error) {
	if id == "" {
		return domain.Employee{}, domain.ErrEmptyID
	}

	emp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			return domain.Employee{}, err
		}
		return domain.Employee{}, fmt.Errorf("service: failed to get employee by id %s: %w", id, err)
	}
	return emp, nil
}

func (s *employeeServiceImpl) CreateEmployee(ctx context.Context, emp domain.Employee) (domain.Employee, error) {
	// Validasi bisnis: cek duplikasi email
	allEmployees, _ := s.repo.FindAll(ctx)
	for _, e := range allEmployees {
		if e.Email == emp.Email {
			return domain.Employee{}, domain.ErrEmailAlreadyExists
		}
	}

	// Validasi bisnis
	if emp.Name == "" {
		return domain.Employee{}, domain.ErrEmptyName
	}
	if emp.Salary < 0 {
		return domain.Employee{}, domain.ErrNegativeSalary
	}
	if emp.Email == "" {
		return domain.Employee{}, domain.ErrEmptyEmail
	}

	created, err := s.repo.Save(ctx, emp)
	if err != nil {
		return domain.Employee{}, fmt.Errorf("service: failed to create employee: %w", err)
	}
	return created, nil
}

func (s *employeeServiceImpl) UpdateEmployee(ctx context.Context, id string, emp domain.Employee) (domain.Employee, error) {
	if id == "" {
		return domain.Employee{}, domain.ErrEmptyID
	}

	// Validasi bisnis: pastikan employee exists sebelum update
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			return domain.Employee{}, fmt.Errorf("service: employee with id %s not found for update: %w", id, err)
		}
		return domain.Employee{}, fmt.Errorf("service: failed to find employee for update: %w", err)
	}

	// Validasi bisnis
	if emp.Name == "" {
		return domain.Employee{}, domain.ErrEmptyName
	}
	if emp.Salary < 0 {
		return domain.Employee{}, domain.ErrNegativeSalary
	}
	if emp.Email == "" {
		return domain.Employee{}, domain.ErrEmptyEmail
	}

	updated, err := s.repo.Update(ctx, id, emp)
	if err != nil {
		return domain.Employee{}, fmt.Errorf("service: failed to update employee with id %s: %w", id, err)
	}
	return updated, nil
}

func (s *employeeServiceImpl) DeleteEmployee(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrEmptyID
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service: failed to delete employee with id %s: %w", id, err)
	}
	return nil
}
