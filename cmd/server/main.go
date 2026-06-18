// cmd/server/main.go
package main

import (
	"decoupled_contract_services/internal/employee"
	"decoupled_contract_services/pkg/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Repository
	// Untuk contoh ini, kita gunakan implementasi in-memory
	employeeRepo := employee.NewInMemoryEmployeeRepository()

	// Jika Anda ingin menggunakan database sungguhan (misal MySQL), Anda perlu:
	// - Menginisialisasi koneksi DB: db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/employee_db")
	// - Membuat instance repository MySQL: employeeRepo := employee.NewEmployeeRepositoryMySQL(db)
	// - Pastikan error handling untuk koneksi DB dilakukan dengan benar

	// 2. Inisialisasi Service dan inject Repository
	employeeSvc := employee.NewEmployeeService(employeeRepo)

	// 3. Inisialisasi Handler dan inject Service
	employeeHandler := handler.NewEmployeeHandler(employeeSvc)

	// 4. Setup Gin Router
	router := gin.Default()

	// Middleware untuk logging dan recovery Gin sudah termasuk secara default
	// Anda bisa menambahkan middleware lain di sini jika perlu (contoh: CORS, authentication)

	// Setup routes API
	api := router.Group("/api/v1/employees")
	{
		api.GET("", employeeHandler.GetAll)
		api.POST("", employeeHandler.Create)
		api.GET("/:id", employeeHandler.GetByID)
		api.PUT("/:id", employeeHandler.Update)
		api.DELETE("/:id", employeeHandler.Delete)
	}

	// Jalankan server
	port := "8080"
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
