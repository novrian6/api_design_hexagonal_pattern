// internal/repositories/postgres/employee_repo.go
package postgres

import (
	"context"
	"fmt"
	"sync"

	"decoupled_contract_services/internal/domain"
	"decoupled_contract_services/internal/ports"
)

// inMemoryEmployeeRepository adalah implementasi in-memory dari EmployeeRepository (outbound port)
// Ini adalah OUTBOUND ADAPTER — untuk production, ganti dengan implementasi PostgreSQL/MySQL sungguhan
type inMemoryEmployeeRepository struct {
	mu     sync.RWMutex
	data   map[string]domain.Employee
	nextID int
}

// NewInMemoryEmployeeRepository membuat instance baru dari repository in-memory
func NewInMemoryEmployeeRepository() ports.EmployeeRepository {
	return &inMemoryEmployeeRepository{
		data:   make(map[string]domain.Employee),
		nextID: 1,
	}
}

func (r *inMemoryEmployeeRepository) FindAll(ctx context.Context) ([]domain.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	employees := make([]domain.Employee, 0, len(r.data))
	for _, emp := range r.data {
		employees = append(employees, emp)
	}
	return employees, nil
}

func (r *inMemoryEmployeeRepository) FindByID(ctx context.Context, id string) (domain.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	emp, exists := r.data[id]
	if !exists {
		return domain.Employee{}, domain.ErrEmployeeNotFound
	}
	return emp, nil
}

func (r *inMemoryEmployeeRepository) Save(ctx context.Context, emp domain.Employee) (domain.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := fmt.Sprintf("EMP-%03d", r.nextID)
	r.nextID++
	emp.ID = id
	r.data[id] = emp
	return emp, nil
}

func (r *inMemoryEmployeeRepository) Update(ctx context.Context, id string, emp domain.Employee) (domain.Employee, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.data[id]
	if !exists {
		return domain.Employee{}, domain.ErrEmployeeNotFound
	}

	emp.ID = id
	r.data[id] = emp
	return emp, nil
}

func (r *inMemoryEmployeeRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.data[id]
	if !exists {
		return domain.ErrEmployeeNotFound
	}

	delete(r.data, id)
	return nil
}
