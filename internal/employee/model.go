// internal/employee/model.go
package employee

// Employee struct merepresentasikan entitas bisnis
type Employee struct {
	ID       string  `json:"id"`
	Name     string  `json:"name" validate:"required"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary" validate:"gte=0"`
	Email    string  `json:"email" validate:"required,email"`
}
