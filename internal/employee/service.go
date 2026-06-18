// internal/employee/service.go
package employee

import (
	"context"
	"fmt"
)

// EmployeeService mendefinisikan kontrak untuk layanan bisnis Employee
type EmployeeService interface {
	GetAllEmployees(ctx context.Context) ([]Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (Employee, error)
	CreateEmployee(ctx context.Context, emp Employee) (Employee, error)
	UpdateEmployee(ctx context.Context, id string, emp Employee) (Employee, error)
	DeleteEmployee(ctx context.Context, id string) error
}

// employeeServiceImpl adalah implementasi konkret dari EmployeeService
type employeeServiceImpl struct {
	repo EmployeeRepository // Bergantung pada interface repository
}

// NewEmployeeService adalah constructor dengan dependency injection
func NewEmployeeService(repo EmployeeRepository) EmployeeService {
	return &employeeServiceImpl{repo: repo}
}

func (s *employeeServiceImpl) GetAllEmployees(ctx context.Context) ([]Employee, error) {
	employees, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get all employees: %w", err)
	}
	// Logika bisnis tambahan bisa ditambahkan di sini, contoh: sorting
	return employees, nil
}

func (s *employeeServiceImpl) GetEmployeeByID(ctx context.Context, id string) (Employee, error) {
	if id == "" {
		return Employee{}, fmt.Errorf("service: employee id cannot be empty")
	}

	emp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Employee{}, fmt.Errorf("service: failed to get employee by id %s: %w", id, err)
	}
	return emp, nil
}

func (s *employeeServiceImpl) CreateEmployee(ctx context.Context, emp Employee) (Employee, error) {
	// Validasi bisnis: cek duplikasi email
	allEmployees, _ := s.repo.FindAll(ctx) // Error handling omitted for brevity in this example
	for _, e := range allEmployees {
		if e.Email == emp.Email {
			return Employee{}, fmt.Errorf("service: email %s already exists", emp.Email)
		}
	}

	// Validasi bisnis tambahan (bisa juga dilakukan di handler menggunakan validator)
	if emp.Name == "" {
		return Employee{}, fmt.Errorf("service: employee name cannot be empty")
	}
	if emp.Salary < 0 {
		return Employee{}, fmt.Errorf("service: salary cannot be negative")
	}
	if emp.Email == "" {
		return Employee{}, fmt.Errorf("service: employee email cannot be empty")
	}

	created, err := s.repo.Save(ctx, emp)
	if err != nil {
		return Employee{}, fmt.Errorf("service: failed to create employee: %w", err)
	}
	return created, nil
}

func (s *employeeServiceImpl) UpdateEmployee(ctx context.Context, id string, emp Employee) (Employee, error) {
	if id == "" {
		return Employee{}, fmt.Errorf("service: employee id cannot be empty")
	}

	// Validasi bisnis: pastikan employee exists sebelum update
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Employee{}, fmt.Errorf("service: employee with id %s not found for update: %w", id, err)
	}

	// Validasi bisnis tambahan (bisa juga dilakukan di handler)
	if emp.Name == "" {
		return Employee{}, fmt.Errorf("service: employee name cannot be empty")
	}
	if emp.Salary < 0 {
		return Employee{}, fmt.Errorf("service: salary cannot be negative")
	}
	if emp.Email == "" {
		return Employee{}, fmt.Errorf("service: employee email cannot be empty")
	}

	updated, err := s.repo.Update(ctx, id, emp)
	if err != nil {
		return Employee{}, fmt.Errorf("service: failed to update employee with id %s: %w", id, err)
	}
	return updated, nil
}

func (s *employeeServiceImpl) DeleteEmployee(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("service: employee id cannot be empty")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service: failed to delete employee with id %s: %w", id, err)
	}
	return nil
}
