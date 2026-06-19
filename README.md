# Employee Service — Hexagonal Architecture (Ports & Adapters)

## 📋 Overview

Project ini adalah implementasi **Hexagonal Architecture** (Ports & Adapters Pattern) untuk REST API manajemen Employee. Struktur ini memisahkan domain core dari infrastruktur eksternal (framework HTTP, database) sehingga aplikasi menjadi lebih **testable**, **maintainable**, dan **swappable**.

## 🏗️ Project Structure

```
employee-service/
├── api/
│   └── openapi.yaml                  # API Specification (OpenAPI 3.0)
├── cmd/
│   └── server/
│       └── main.go                   # Composition Root — Wiring dependency
├── config/
│   └── config.yaml                   # Konfigurasi aplikasi
├── internal/
│   ├── domain/                       # INTI HEXAGON — Tanpa dependensi framework
│   │   ├── employee.go               # Entity: type Employee struct
│   │   └── errors.go                 # Sentinel errors untuk domain
│   ├── ports/                        # PORT — Hanya berisi interface
│   │   ├── inbound.go                # EmployeeService interface
│   │   └── outbound.go               # EmployeeRepository interface
│   ├── services/                     # USE CASE — Implementasi inbound port
│   │   └── employee_service.go       # Business logic dengan dependency ke outbound port
│   ├── handlers/                     # INBOUND ADAPTER — HTTP → Domain
│   │   └── rest/
│   │       └── employee_handler.go   # Gin handler, memanggil EmployeeService
│   ├── repositories/                 # OUTBOUND ADAPTER — Domain → Database
│   │   └── postgres/
│   │       └── employee_repo.go      # Implementasi EmployeeRepository (in-memory)
│   └── middleware/                   # Cross-cutting concerns
│       └── auth.go                   # Auth, CORS, Logging middleware
├── pkg/
│   └── response/                     # Helper untuk response JSON standar
├── scripts/
│   └── build.sh                      # Build script
└── test/
    └── fixtures/
        └── employee_testdata.json    # Test data
```

## 🔄 Dependency Flow

```
[Inbound Adapter] ──depends on──▶ [Inbound Port] (Interface)
  (HTTP Handler)                          │
      Gin                                 ▼
                                   [Domain Logic] (Implementation)
                                     (Service Layer)
                                          │
                                          ▼
[Outbound Adapter] ◀──implements─── [Outbound Port] (Interface)
    (Repository)                       (Repository Interface)
```

## 🚀 Quick Start

### Prerequisites
- Go 1.25+
- (Optional) `jq` untuk formatting JSON output

### Run the server

```bash
# Build & run
go run ./cmd/server/

# Atau build binary
go build -o bin/employee-server ./cmd/server/
./bin/employee-server
```

Server akan berjalan di `http://localhost:8080`

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/employees` | Get all employees |
| POST | `/api/v1/employees` | Create employee |
| GET | `/api/v1/employees/:id` | Get employee by ID |
| PUT | `/api/v1/employees/:id` | Update employee |
| DELETE | `/api/v1/employees/:id` | Delete employee |

### Example curl commands

```bash
# Health check
curl -s http://localhost:8080/health | jq .

# Get all employees
curl -s http://localhost:8080/api/v1/employees | jq .

# Create employee
curl -s -X POST http://localhost:8080/api/v1/employees \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","position":"Software Engineer","salary":5000000,"email":"john@example.com"}' | jq .

# Get by ID
curl -s http://localhost:8080/api/v1/employees/EMP-001 | jq .

# Update
curl -s -X PUT http://localhost:8080/api/v1/employees/EMP-001 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated","position":"Senior Engineer","salary":7000000,"email":"john@example.com"}' | jq .

# Delete
curl -s -X DELETE http://localhost:8080/api/v1/employees/EMP-001
```

## 🧪 Testing

### Manual test with curl
```bash
# Run all test scenarios
chmod +x scripts/build.sh
./scripts/build.sh
./bin/employee-server &
# Then run curl commands above
```

## 📚 Key Concepts

| Layer | Description | Dependencies |
|-------|-------------|--------------|
| **Domain** | Pure business logic & entities | None (zero dependencies) |
| **Ports** | Interface contracts | Only domain |
| **Services** | Use case implementation | Domain + Ports |
| **Handlers** | HTTP adapters | Ports (interface only) |
| **Repositories** | Database adapters | Domain + Ports |

## 📖 Documentation

- [HEXAGONAL_PATTERN.md](./HEXAGONAL_PATTERN.md) — Detailed analysis of architecture
- [api/openapi.yaml](./api/openapi.yaml) — API specification# api_design_hexagonal_pattern_pure
