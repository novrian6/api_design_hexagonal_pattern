// internal/ports/outbound.go
package ports

import (
	"context"

	"decoupled_contract_services/internal/domain"
)

// EmployeeRepository mendefinisikan kontrak untuk akses data Employee
// Ini adalah OUTBOUND PORT — milik domain core, diimplementasikan oleh adapter (repository)
type EmployeeRepository interface {
	FindAll(ctx context.Context) ([]domain.Employee, error)
	FindByID(ctx context.Context, id string) (domain.Employee, error)
	Save(ctx context.Context, emp domain.Employee) (domain.Employee, error)
	Update(ctx context.Context, id string, emp domain.Employee) (domain.Employee, error)
	Delete(ctx context.Context, id string) error
}
