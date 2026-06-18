// pkg/handler/employee_handler.go
package handler

import (
	"fmt"
	"net/http"
	"decoupled_contract_services/internal/employee"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10" // Untuk validasi input
)

type EmployeeHandler struct {
	service  employee.EmployeeService // Bergantung pada interface service
	validate *validator.Validate
}

func NewEmployeeHandler(svc employee.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service:  svc,
		validate: validator.New(), // Inisialisasi validator
	}
}

// GetAll handles GET /api/v1/employees
func (h *EmployeeHandler) GetAll(c *gin.Context) {
	employees, err := h.service.GetAllEmployees(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

// GetByID handles GET /api/v1/employees/:id
func (h *EmployeeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	emp, err := h.service.GetEmployeeByID(c.Request.Context(), id)
	if err != nil {
		// Periksa apakah error spesifik bahwa ID tidak ditemukan
		if err.Error() == fmt.Sprintf("service: employee with id %s not found", id) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, emp)
}

// Create handles POST /api/v1/employees
func (h *EmployeeHandler) Create(c *gin.Context) {
	var emp employee.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gunakan validator untuk validasi struct
	if err := h.validate.Struct(emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.CreateEmployee(c.Request.Context(), emp)
	if err != nil {
		// Tangani error spesifik dari service, contoh duplikasi email
		if err.Error() == fmt.Sprintf("service: email %s already exists", emp.Email) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, created)
}

// Update handles PUT /api/v1/employees/:id
func (h *EmployeeHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var emp employee.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gunakan validator untuk validasi struct
	if err := h.validate.Struct(emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.service.UpdateEmployee(c.Request.Context(), id, emp)
	if err != nil {
		// Tangani error spesifik
		if err.Error() == fmt.Sprintf("service: employee with id %s not found for update", id) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, updated)
}

// Delete handles DELETE /api/v1/employees/:id
func (h *EmployeeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteEmployee(c.Request.Context(), id)
	if err != nil {
		if err.Error() == fmt.Sprintf("service: failed to delete employee with id %s: employee with id %s not found", id, id) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent) // 204 No Content
}
