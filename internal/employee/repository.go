// internal/employee/repository.go
package employee

import "context"

// EmployeeRepository mendefinisikan kontrak untuk akses data Employee
type EmployeeRepository interface {
	FindAll(ctx context.Context) ([]Employee, error)
	FindByID(ctx context.Context, id string) (Employee, error)
	Save(ctx context.Context, emp Employee) (Employee, error)
	Update(ctx context.Context, id string, emp Employee) (Employee, error)
	Delete(ctx context.Context, id string) error
}
