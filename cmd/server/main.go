// cmd/server/main.go
package main

import (
	"log"

	"decoupled_contract_services/internal/handlers/rest"
	"decoupled_contract_services/internal/middleware"
	postgresrepo "decoupled_contract_services/internal/repositories/postgres"
	"decoupled_contract_services/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Repository (Outbound Adapter)
	// Untuk contoh ini, kita gunakan implementasi in-memory
	employeeRepo := postgresrepo.NewInMemoryEmployeeRepository()

	// Jika Anda ingin menggunakan database sungguhan (misal PostgreSQL), Anda perlu:
	// - Menginisialisasi koneksi DB: db, err := sql.Open("postgres", "user:password@localhost:5432/employee_db?sslmode=disable")
	// - Membuat instance repository PostgreSQL: employeeRepo := postgresrepo.NewEmployeeRepositoryPostgres(db)
	// - Pastikan error handling untuk koneksi DB dilakukan dengan benar

	// 2. Inisialisasi Service (Use Case) dan inject Repository
	employeeSvc := services.NewEmployeeService(employeeRepo)

	// 3. Inisialisasi Handler (Inbound Adapter) dan inject Service
	employeeHandler := rest.NewEmployeeHandler(employeeSvc)

	// 4. Setup Gin Router
	router := gin.Default()

	// Middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

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
	log.Printf("API available at http://localhost:%s/api/v1/employees", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
