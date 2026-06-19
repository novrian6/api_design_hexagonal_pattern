// internal/handlers/rest/employee_handler.go
package rest

import (
	"errors"
	"net/http"

	"decoupled_contract_services/internal/domain"
	"decoupled_contract_services/internal/ports"
	"decoupled_contract_services/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// EmployeeHandler adalah INBOUND ADAPTER untuk HTTP (REST API)
// Bergantung pada interface EmployeeService (inbound port), bukan implementasi konkret
type EmployeeHandler struct {
	service  ports.EmployeeService
	validate *validator.Validate
}

// NewEmployeeHandler adalah constructor dengan dependency injection
func NewEmployeeHandler(svc ports.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service:  svc,
		validate: validator.New(),
	}
}

// GetAll handles GET /api/v1/employees
func (h *EmployeeHandler) GetAll(c *gin.Context) {
	employees, err := h.service.GetAllEmployees(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, employees)
}

// GetByID handles GET /api/v1/employees/:id
func (h *EmployeeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	emp, err := h.service.GetEmployeeByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, emp)
}

// Create handles POST /api/v1/employees
func (h *EmployeeHandler) Create(c *gin.Context) {
	var emp domain.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Gunakan validator untuk validasi struct
	if err := h.validate.Struct(emp); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.service.CreateEmployee(c.Request.Context(), emp)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, created)
}

// Update handles PUT /api/v1/employees/:id
func (h *EmployeeHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var emp domain.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Gunakan validator untuk validasi struct
	if err := h.validate.Struct(emp); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.service.UpdateEmployee(c.Request.Context(), id, emp)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, updated)
}

// Delete handles DELETE /api/v1/employees/:id
func (h *EmployeeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteEmployee(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrEmployeeNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent) // 204 No Content
}
