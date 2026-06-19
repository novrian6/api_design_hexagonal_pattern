// pkg/response/response.go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response adalah struktur standar untuk response JSON API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success mengirim response sukses
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMessage mengirim response sukses dengan pesan
func SuccessWithMessage(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error mengirim response error
func Error(c *gin.Context, statusCode int, errMsg string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   errMsg,
	})
}

// Paginated digunakan untuk response dengan pagination
type Paginated struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalItems int         `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

// SuccessPaginated mengirim response sukses dengan pagination
func SuccessPaginated(c *gin.Context, page, perPage, totalItems int, data interface{}) {
	totalPages := totalItems / perPage
	if totalItems%perPage != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: Paginated{
			Data:       data,
			Page:       page,
			PerPage:    perPage,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	})
}
