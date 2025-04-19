# Panduan Pengembangan API KijunPOS Go

Dokumen ini berisi panduan dan aturan untuk mengembangkan API baru pada project KijunPOS Go menggunakan prinsip Clean Architecture dan pendekatan modular.

## Struktur Project

```
kijunpos-go/
├── cmd/                    # Command-line entry points
├── config/                 # Konfigurasi aplikasi
├── internal/               # Kode private aplikasi
│   ├── app/                # Application bootstrapping dan DI
│   ├── domain/             # Domain entities dan interfaces
│   ├── usecase/            # Business logic
│   │   └── [domain]/       # Usecase per domain bisnis
│   ├── repository/         # Data access implementations
│   │   └── [domain]/       # Repository per domain bisnis
│   ├── delivery/           # Delivery mechanisms (gRPC, HTTP, etc)
│   │   └── grpc/           # gRPC handlers
│   │       └── [domain]/   # Handler per domain bisnis
│   └── pkg/                # Internal shared packages
├── pkg/                    # Public shared packages (jika ada)
└── proto/                  # Proto definitions
```

## Langkah-langkah Membuat API Baru

Berikut adalah langkah-langkah untuk membuat API baru dengan mengikuti prinsip Clean Architecture dan pendekatan modular:

### 1. Mendefinisikan Proto (Contract)

1. Buat atau update file proto di direktori `proto/[domain]/[domain].proto`
2. Definisikan service, request, dan response messages
3. Generate kode Go dari proto dengan menjalankan `buf generate`

Contoh:
```protobuf
// proto/product/product.proto
syntax = "proto3";

package product;

option go_package = "./product";

service ProductService {
  rpc CreateProduct(CreateProductRequest) returns (ProductResponse) {}
  rpc GetProduct(GetProductRequest) returns (ProductResponse) {}
}

message CreateProductRequest {
  string name = 1;
  string description = 2;
  double price = 3;
}

message GetProductRequest {
  string id = 1;
}

message ProductResponse {
  string id = 1;
  string name = 2;
  string description = 3;
  double price = 4;
  bool success = 5;
  string message = 6;
}
```

### 2. Mendefinisikan Domain Entity dan Interface

1. Buat file `internal/domain/[domain].go`
2. Definisikan struct entity
3. Definisikan interface repository
4. Definisikan interface usecase

Contoh:
```go
// internal/domain/product.go
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Product represents the product entity
type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

// ProductRepository represents the product repository contract
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ProductUseCase represents the product use case contract
type ProductUseCase interface {
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}
```

### 3. Implementasi Repository

1. Buat direktori `internal/repository/[domain]/`
2. Buat file `internal/repository/[domain]/[domain].go` untuk struct dan konstruktor
3. Buat file terpisah untuk setiap method repository
4. Update `internal/repository/factory.go` untuk menambahkan factory function

Contoh:
```go
// internal/repository/product/product.go
package product

import (
	"github/kijunpos/config/db"
	"github/kijunpos/internal/domain"
)

type productRepository struct {
	dbConn *db.Connection
}

// NewProductRepository creates a new product repository
func NewProductRepository(dbConn *db.Connection) domain.ProductRepository {
	return &productRepository{
		dbConn: dbConn,
	}
}
```

```go
// internal/repository/product/create.go
package product

import (
	"context"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
)

// Create creates a new product in the database
func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	ctx, span := apm.GetTracer().Start(ctx, "repository.product.Create")
	defer span.End()

	query := `
		INSERT INTO products (
			id, name, description, price, created_at
		) VALUES (
			$1, $2, $3, $4, $5
		)
	`

	_, err := r.dbConn.DB.ExecContext(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.CreatedAt,
	)

	return err
}
```

```go
// internal/repository/factory.go
// Tambahkan factory function untuk product repository
func NewProductRepository(dbConn *db.Connection) domain.ProductRepository {
	return productRepo.NewProductRepository(dbConn)
}
```

### 4. Implementasi Usecase

1. Buat direktori `internal/usecase/[domain]/`
2. Buat file `internal/usecase/[domain]/[domain].go` untuk struct dan konstruktor
3. Buat file terpisah untuk setiap method usecase

Contoh:
```go
// internal/usecase/product/product.go
package product

import (
	"github/kijunpos/internal/domain"
)

type productUseCase struct {
	productRepo domain.ProductRepository
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(productRepo domain.ProductRepository) domain.ProductUseCase {
	return &productUseCase{
		productRepo: productRepo,
	}
}
```

```go
// internal/usecase/product/create_product.go
package product

import (
	"context"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
	"time"

	"github.com/google/uuid"
)

// CreateProduct creates a new product
func (uc *productUseCase) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	ctx, span := apm.GetTracer().Start(ctx, "usecase.product.CreateProduct")
	defer span.End()

	// Prepare product data
	product.ID = uuid.New()
	product.CreatedAt = time.Now()

	// Create product in the database
	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}
```

### 5. Implementasi Delivery (Handler)

1. Buat direktori `internal/delivery/grpc/[domain]/`
2. Buat file `internal/delivery/grpc/[domain]/handler.go` untuk struct dan konstruktor
3. Buat file terpisah untuk setiap method handler
4. Update `internal/delivery/grpc/factory.go` untuk menambahkan factory function

Contoh:
```go
// internal/delivery/grpc/product/handler.go
package product

import (
	"github/kijunpos/internal/domain"
	pbProduct "github/kijunpos/gen/proto/product"
)

// Handler handles gRPC requests for product service
type Handler struct {
	pbProduct.UnimplementedProductServiceServer
	productUseCase domain.ProductUseCase
}

// NewHandler creates a new product handler
func NewHandler(productUseCase domain.ProductUseCase) *Handler {
	return &Handler{
		productUseCase: productUseCase,
	}
}
```

```go
// internal/delivery/grpc/product/create_product.go
package product

import (
	"context"
	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/pkg/apm"
	pbProduct "github/kijunpos/gen/proto/product"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateProduct handles product creation
func (h *Handler) CreateProduct(ctx context.Context, req *pbProduct.CreateProductRequest) (*pbProduct.ProductResponse, error) {
	ctx, span := apm.GetTracer().Start(ctx, "delivery.grpc.product.CreateProduct")
	defer span.End()

	// Validate request
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "product name is required")
	}

	// Create product domain object
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	// Call use case
	createdProduct, err := h.productUseCase.CreateProduct(ctx, product)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbProduct.ProductResponse{
		Id:          createdProduct.ID.String(),
		Name:        createdProduct.Name,
		Description: createdProduct.Description,
		Price:       createdProduct.Price,
		Success:     true,
		Message:     "Product created successfully",
	}, nil
}
```

```go
// internal/delivery/grpc/factory.go
// Tambahkan interface dan factory function untuk product handler
type ProductServiceHandler interface {
	pbProduct.ProductServiceServer
}

func NewProductServiceHandler(productUseCase domain.ProductUseCase) ProductServiceHandler {
	return productHandler.NewHandler(productUseCase)
}
```

### 6. Update Server untuk Registrasi Service

1. Update `internal/delivery/grpc/server.go` untuk menambahkan registrasi service baru

```go
// internal/delivery/grpc/server.go
// Tambahkan parameter productHandler
func StartGRPCServer(cfg *config.Config, userHandler UserServiceHandler, productHandler ProductServiceHandler) {
	// ...

	// Register services
	pbUser.RegisterUserServiceServer(grpcServer, userHandler)
	pbProduct.RegisterProductServiceServer(grpcServer, productHandler)

	// ...
}
```

### 7. Update Application untuk Dependency Injection

1. Update `internal/app/app.go` untuk menambahkan dependency injection

```go
// internal/app/app.go
// Tambahkan field untuk ProductHandler
type Application struct {
	Config          *config.Config
	DBManager       *db.Manager
	UserHandler     grpc.UserServiceHandler
	ProductHandler  grpc.ProductServiceHandler
}

// Update NewApplication untuk inisialisasi product repository, usecase, dan handler
func NewApplication() *Application {
	// ...

	// Initialize repositories
	userRepo := repository.NewUserRepository(kijunConn)
	productRepo := repository.NewProductRepository(kijunConn)

	// Initialize use cases
	userUC := userUseCase.NewUserUseCase(userRepo)
	productUC := productUseCase.NewProductUseCase(productRepo)

	// Initialize gRPC handlers
	userHandler := grpc.NewUserServiceHandler(userUC)
	productHandler := grpc.NewProductServiceHandler(productUC)

	return &Application{
		Config:          configData,
		DBManager:       dbManager,
		UserHandler:     userHandler,
		ProductHandler:  productHandler,
	}
}

// Update Start untuk menambahkan parameter productHandler
func (app *Application) Start() {
	// Start the gRPC server
	grpc.StartGRPCServer(app.Config, app.UserHandler, app.ProductHandler)
}
```

## Aturan dan Konvensi

### Penamaan

1. **File**: Gunakan snake_case untuk nama file (contoh: `create_product.go`)
2. **Package**: Gunakan lowercase tanpa underscore (contoh: `package product`)
3. **Struct**: Gunakan PascalCase untuk struct yang di-export (contoh: `type Product struct`)
4. **Interface**: Gunakan PascalCase dan akhiri dengan kata yang menjelaskan perannya (contoh: `type ProductRepository interface`)
5. **Method**: Gunakan PascalCase untuk method yang di-export (contoh: `func (r *productRepository) Create()`)
6. **Variable**: Gunakan camelCase untuk variable (contoh: `productRepo`)

### Struktur File

1. **Domain**: Satu file per domain, berisi entity dan interface
2. **Repository**: Satu file per method repository
3. **Usecase**: Satu file per method usecase
4. **Delivery**: Satu file per method handler

### Dependency Injection

1. Gunakan constructor pattern untuk dependency injection
2. Gunakan factory pattern untuk membuat instance
3. Injeksi dependency dari luar, bukan dibuat di dalam komponen

### Error Handling

1. Gunakan error wrapping untuk menambahkan konteks pada error
2. Gunakan custom error untuk error domain
3. Jangan expose error internal ke client, konversi ke error yang sesuai

### Logging dan Tracing

1. Gunakan apm untuk tracing
2. Gunakan logger untuk logging
3. Tambahkan span untuk setiap method

## Testing

### Unit Testing

1. Buat file test untuk setiap file implementasi
2. Gunakan mock untuk dependency
3. Fokus pada satu unit functionality

Contoh:
```go
// internal/usecase/product/create_product_test.go
package product_test

import (
	"context"
	"testing"
	"time"

	"github/kijunpos/internal/domain"
	"github/kijunpos/internal/domain/mocks"
	"github/kijunpos/internal/usecase/product"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduct(t *testing.T) {
	mockRepo := new(mocks.ProductRepository)
	useCase := product.NewProductUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		productData := &domain.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       100.0,
		}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil)

		result, err := useCase.CreateProduct(context.Background(), productData)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEqual(t, uuid.Nil, result.ID)
		assert.WithinDuration(t, time.Now(), result.CreatedAt, 2*time.Second)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		productData := &domain.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       100.0,
		}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(errors.New("repository error"))

		result, err := useCase.CreateProduct(context.Background(), productData)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
```

### Integration Testing

1. Buat test yang menguji interaksi antar komponen
2. Gunakan database test

## Kesimpulan

Dengan mengikuti panduan ini, pengembangan API baru pada project KijunPOS Go akan konsisten, terstruktur, dan mudah dipelihara. Pendekatan modular memudahkan pengujian dan pemeliharaan kode, sementara Clean Architecture memastikan pemisahan tanggung jawab yang jelas antar layer.