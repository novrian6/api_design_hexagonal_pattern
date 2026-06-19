// internal/ports/inbound.go
package ports

import (
	"context"

	"decoupled_contract_services/internal/domain"
)

// EmployeeService mendefinisikan kontrak untuk use case / layanan bisnis Employee
// Ini adalah INBOUND PORT — milik domain core, bebas dari framework
type EmployeeService interface {
	GetAllEmployees(ctx context.Context) ([]domain.Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (domain.Employee, error)
	CreateEmployee(ctx context.Context, emp domain.Employee) (domain.Employee, error)
	UpdateEmployee(ctx context.Context, id string, emp domain.Employee) (domain.Employee, error)
	DeleteEmployee(ctx context.Context, id string) error
}
