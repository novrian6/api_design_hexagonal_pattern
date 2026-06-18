# Hexagonal Architecture Analysis

> Project: `decoupled_contract_services`  
> Analisis struktur dan implementasi Hexagonal Architecture (Ports & Adapters Pattern)

---

## Table of Contents

1. [Apa itu Hexagonal Architecture?](#apa-itu-hexagonal-architecture)
2. [Struktur Project](#struktur-project)
3. [Analisis Per Layer](#analisis-per-layer)
4. [Diagram Dependency Flow](#diagram-dependency-flow)
5. [Kesesuaian dengan Prinsip Hexagonal](#kesesuaian-dengan-prinsip-hexagonal)
6. [Improvement Opportunities](#improvement-opportunities)
7. [Benchmark: Sebelum vs Sesudah Hexagonal](#benchmark-sebelum-vs-sesudah-hexagonal)
8. [Kesimpulan](#kesimpulan)

---

## Apa itu Hexagonal Architecture?

**Hexagonal Architecture** (juga dikenal sebagai **Ports & Adapters Pattern**) diperkenalkan oleh Alistair Cockburn. Tujuan utamanya adalah menciptakan aplikasi yang:

- **Terisolasi** dari teknologi eksternal (framework, database, UI)
- **Testable** — domain logic bisa diuji tanpa infrastruktur
- **Swappable** — komponen infrastruktur bisa diganti tanpa mengubah logika bisnis

### Konsep Inti

```
                    ┌──────────────────────────────┐
                    │        Application Core       │
                    │   ┌──────────────────────┐   │
  ┌──────────┐      │   │    Domain Entity     │   │      ┌──────────────┐
  │  HTTP    │──────┼──▶│                      │◀──┼──────│   Database   │
  │ Adapter  │      │   │    Use Case (Port)   │   │      │   Adapter    │
  │(Inbound) │◀────┼──▶│    Service (Impl)     │───┼──────│ (Outbound)   │
  └──────────┘      │   │                      │   │      └──────────────┘
                    │   │    Repository (Port) │   │
                    │   └──────────────────────┘   │
                    └──────────────────────────────┘
```

### Aliran Dependency

```
[Inbound Adapter] ──depends on──▶ [Inbound Port] (Interface)
                                          │
                                          ▼
                                   [Domain Logic] (Implementation)
                                          │
                                          ▼
[Outbound Adapter] ◀──implements─── [Outbound Port] (Interface)
```

**Aturan Utama:** Dependency hanya boleh mengarah ke **dalam** (menuju domain core). Layer luar hanya boleh bergantung pada abstraksi di layer dalam.

---

## Struktur Project

```
hexagonal_pattern/
├── cmd/
│   └── server/
│       └── main.go              ← Composition Root (wiring semua dependency)
├── config/
│   └── config.go                ← Konfigurasi aplikasi
├── internal/
│   └── employee/
│       ├── model.go             ← Domain Entity (core bisnis)
│       ├── repository.go        ← Outbound Port (interface)
│       ├── repository_impl.go   ← Outbound Adapter (implementasi)
│       └── service.go           ← Inbound Port + Implementasi
├── pkg/
│   └── handler/
│       └── employee_handler.go  ← Inbound Adapter (HTTP handler)
├── go.mod
└── go.sum
```

---

## Analisis Per Layer

### 1. Domain Entity — `internal/employee/model.go`

```go
type Employee struct {
    ID       string  `json:"id"`
    Name     string  `json:"name" validate:"required"`
    Position string  `json:"position"`
    Salary   float64 `json:"salary" validate:"gte=0"`
    Email    string  `json:"email" validate:"required,email"`
}
```

| Kriteria | Status | Keterangan |
|----------|--------|------------|
| Murni domain, tanpa framework dependency | ✅ | Struct Go murni, hanya tag validasi |
| Representasi bisnis | ✅ | Employee sebagai entitas inti |
| Bisa di-test tanpa infrastruktur | ✅ | Tidak bergantung pada DB/HTTP |

### 2. Inbound Port — `internal/employee/service.go` (Interface)

```go
type EmployeeService interface {
    GetAllEmployees(ctx context.Context) ([]Employee, error)
    GetEmployeeByID(ctx context.Context, id string) (Employee, error)
    CreateEmployee(ctx context.Context, emp Employee) (Employee, error)
    UpdateEmployee(ctx context.Context, id string, emp Employee) (Employee, error)
    DeleteEmployee(ctx context.Context, id string) error
}
```

| Kriteria | Status | Keterangan |
|----------|--------|------------|
| Mendefinisikan kontrak use case | ✅ | Semua operasi bisnis Employee |
| Interface di sisi domain | ✅ | Berada di `internal/employee/` |
| Bisa diimplementasikan berbagai adapter | ✅ | HTTP, gRPC, CLI, dll |

### 3. Inbound Adapter — `pkg/handler/employee_handler.go`

```go
type EmployeeHandler struct {
    service  employee.EmployeeService  // Bergantung pada INTERFACE, bukan implementasi
    validate *validator.Validate
}
```

| Kriteria | Status | Keterangan |
|----------|--------|------------|
| Bergantung pada interface | ✅ | `EmployeeService` interface |
| Menerjemahkan input eksternal ke domain | ✅ | HTTP request → service call → HTTP response |
| Bisa diganti adapter lain | ✅ | Bisa ditambah gRPC handler tanpa modifikasi service |

### 4. Outbound Port — `internal/employee/repository.go`

```go
type EmployeeRepository interface {
    FindAll(ctx context.Context) ([]Employee, error)
    FindByID(ctx context.Context, id string) (Employee, error)
    Save(ctx context.Context, emp Employee) (Employee, error)
    Update(ctx context.Context, id string, emp Employee) (Employee, error)
    Delete(ctx context.Context, id string) error
}
```

| Kriteria | Status | Keterangan |
|----------|--------|------------|
| Kontrak akses data | ✅ | CRUD operations |
| Interface di domain | ✅ | Berada di `internal/employee/` |
| Bisa diimplementasikan berbagai database | ✅ | In-memory, MySQL, PostgreSQL, MongoDB |

### 5. Outbound Adapter — `internal/employee/repository_impl.go`

```go
type inMemoryEmployeeRepository struct {
    mu     sync.RWMutex
    data   map[string]Employee
    nextID int
}
```

| Kriteria | Status | Keterangan |
|----------|--------|------------|
| Mengimplementasikan interface | ✅ | `EmployeeRepository` |
| Bisa diganti tanpa efek ke domain | ✅ | Service tidak tahu implementasi |
| Testable | ✅ | Bisa di-test dengan mock/stub |

### 6. Composition Root — `cmd/server/main.go`

```go
func main() {
    // 1. Init Repository (Outbound Adapter)
    employeeRepo := employee.NewInMemoryEmployeeRepository()

    // 2. Init Service + inject Repository
    employeeSvc := employee.NewEmployeeService(employeeRepo)

    // 3. Init Handler + inject Service
    employeeHandler := handler.NewEmployeeHandler(employeeSvc)

    // 4. Setup router
    router := gin.Default()
    api := router.Group("/api/v1/employees")
    {
        api.GET("", employeeHandler.GetAll)
        api.POST("", employeeHandler.Create)
        api.GET("/:id", employeeHandler.GetByID)
        api.PUT("/:id", employeeHandler.Update)
        api.DELETE("/:id", employeeHandler.Delete)
    }

    router.Run(":8080")
}
```

| Kriteria | Status | Keterangan |
|----------|--------|------------|
| Wiring dependency di satu tempat | ✅ | Semua DI dilakukan di main.go |
| Tidak ada hardcode dependency di layer lain | ✅ | Tidak ada `new` di handler/service |

---

## Diagram Dependency Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                         cmd/server/main.go                          │
│                  (Composition Root — Wiring DI)                     │
└────────────────────┬──────────────────────┬─────────────────────────┘
                     │                      │
                     ▼                      ▼
┌──────────────────────────┐   ┌──────────────────────────────────────┐
│   INBOUND ADAPTER        │   │   OUTBOUND PORT + ADAPTER            │
│   pkg/handler/           │   │   internal/employee/                 │
│   employee_handler.go    │   │                                      │
│                          │   │   ┌──────────────────────────────┐   │
│   Bergantung pada:       │   │   │  EmployeeRepository (Port)   │   │
│   EmployeeService (iface)│   │   │  └ interface{ FindAll... }   │   │
│                          │   │   └──────────────────────────────┘   │
└──────────┬───────────────┘   │               ▲                      │
           │                   │               │                      │
           │ depends on        │     implements│                      │
           ▼                   │               │                      │
┌──────────────────────────┐   │   ┌──────────────────────────────┐   │
│   INBOUND PORT + CORE    │   │   │  InMemoryEmployeeRepository  │   │
│   internal/employee/     │   │   │  (Outbound Adapter)          │   │
│   service.go             │   │   │  └ implements EmployeeRepo   │   │
│                          │   │   └──────────────────────────────┘   │
│   EmployeeService (iface)│   │                                      │
│   employeeServiceImpl    │   │   ┌──────────────────────────────┐   │
│                          │   │   │  EmployeeRepository (Port)   │   │
│   Bergantung pada:       │   │   │  └ interface{ Save... }      │   │
│   EmployeeRepository     │   │   └──────────────────────────────┘   │
└──────────────────────────┘   └──────────────────────────────────────┘
           │
           ▼
┌──────────────────────────────┐
│   DOMAIN ENTITY              │
│   internal/employee/         │
│   model.go                   │
│                              │
│   Employee (pure struct)     │
└──────────────────────────────┘
```

---

## Kesesuaian dengan Prinsip Hexagonal

| Prinsip | Status | Implementasi |
|---------|--------|--------------|
| **Core Domain Terisolasi** | ✅ | `model.go` dan `service.go` tidak import framework eksternal |
| **Dependency Inversion** | ✅ | Semua layer bergantung pada interface (abstraksi), bukan konkret |
| **Dependency Injection** | ✅ | Wiring di `main.go` — composition root |
| **Ports (Interfaces) di Domain** | ✅ | `EmployeeService` & `EmployeeRepository` di `internal/employee/` |
| **Adapters di Layer Luar** | ✅ | Handler di `pkg/handler/`, impl repository di `internal/employee/` |
| **Separation of Concerns** | ✅ | Handler → Service → Repository, masing-masing punya 1 tanggung jawab |
| **Testability** | ✅ | Service bisa di-test dengan mock repository |
| **Swappable Infrastructure** | ✅ | Bisa ganti HTTP → gRPC, InMemory → MySQL tanpa ubah service |
| **Inbound/Outbound Separation** | ✅ | Interface inbound (service) dan outbound (repository) terpisah |

---

## Improvement Opportunities

### 1. Package Structure — Pisahkan Outbound Adapter dari Domain

**Current:**
```
internal/employee/
├── model.go              ← Domain
├── repository.go         ← Port (domain)
├── repository_impl.go    ← Adapter (infrastruktur) — ❌ masih satu package
└── service.go            ← Port + Implementasi
```

**Rekomendasi:**
```
internal/employee/
├── model.go                       ← Domain Entity
├── repository.go                  ← Outbound Port (interface)
├── service.go                     ← Inbound Port + Business Logic
│
internal/employee/repository/
├── memory/
│   └── employee.go               ← In-Memory Adapter
├── mysql/
│   └── employee.go               ← MySQL Adapter
└── postgres/
    └── employee.go               ← PostgreSQL Adapter
```

**Keuntungan:**
- Domain benar-benar bebas dari infrastruktur
- Adapter bisa di-deploy secara independen
- Builder pattern lebih jelas saat menambah adapter baru

### 2. Error Handling — Gunakan Sentinel Errors

**Current (❌ String comparison):**
```go
// di handler
if err.Error() == fmt.Sprintf("service: email %s already exists", emp.Email) {
    c.JSON(http.StatusConflict, ...)
}
```

**Rekomendasi:**
```go
// di service.go (domain layer)
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrEmployeeNotFound = errors.New("employee not found")

// di handler
if errors.Is(err, employee.ErrEmailAlreadyExists) {
    c.JSON(http.StatusConflict, ...)
}
```

**Keuntungan:**
- Decoupling error handling antar layer
- Type-safe dan bisa di-test
- Handler tidak perlu tahu format string error

### 3. Validasi Terpusat

**Current:** Validasi terbagi antara handler (`validator` library) dan service (manual check).
**Rekomendasi:** Semua validasi di service layer, handler hanya untuk parsing input.

### 4. Konfigurasi Terpusat

**Current:** `config/config.go` minimal.
**Rekomendasi:** Environment-based configuration untuk database, port, environment, dll.

---

## Benchmark: Sebelum vs Sesudah Hexagonal

| Aspek | Tanpa Hexagonal (Monolitik) | Dengan Hexagonal ✅ |
|-------|----------------------------|---------------------|
| **Testabilitas** | Sulit — perlu setup DB/HTTP | Mudah — mock interface |
| **Flexibilitas DB** | Tight coupling ke SQL | Ganti adapter tanpa efek |
| **Flexibilitas Framework** | Tight coupling ke Gin/Gorilla | Bisa ganti HTTP framework |
| **Maintainability** | Logic tersebar | Separation of concerns jelas |
| **Onboarding Developer** | Perlu paham semua stack | Bisa fokus per layer |
| **Unit Test Coverage** | Rendah — dependency berat | Tinggi — mudah mocking |
| **Code Reusability** | Rendah — tied ke framework | Tinggi — domain reusable |

---

## Kesimpulan

### ✅ Verdict: **Hexagonal Architecture Terimplementasi dengan Baik**

Project `decoupled_contract_services` sudah menerapkan prinsip-prinsip Hexagonal Architecture secara konsisten:

1. **Ports & Adapters terdefinisi jelas**
   - Inbound Port: `EmployeeService` interface
   - Outbound Port: `EmployeeRepository` interface
   - Inbound Adapter: `EmployeeHandler` (HTTP)
   - Outbound Adapter: `inMemoryEmployeeRepository`

2. **Dependency Inversion Principle terpenuhi**
   - Handler → Service (interface)
   - Service → Repository (interface)
   - Tidak ada ketergantungan ke implementasi konkret

3. **Domain model murni**
   - Struct `Employee` tanpa framework dependency
   - Business logic terpusat di service

4. **Extensible & Testable**
   - Siap untuk implementasi database sungguhan (MySQL, PostgreSQL)
   - Siap untuk tambahan adapter inbound (gRPC, GraphQL, CLI)
   - Setiap layer bisa di-unit test dengan mocking

> **Struktur ini adalah fondasi yang solid untuk pengembangan lebih lanjut.**
> Dengan menambahkan pemisahan package adapter dan sentinel errors, project ini bisa mencapai tingkat kematangan arsitektur yang optimal.

---

## Referensi

- [Alistair Cockburn — Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Ports & Adapters Pattern](https://docs.microsoft.com/en-us/dotnet/architecture/modern-web-apps-azure/common-web-application-architectures#hexagonal-architect)
- [Go Lang — Project Layout Standards](https://github.com/golang-standards/project-layout)