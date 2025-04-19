# Panduan Migrasi ke Clean Architecture

## Struktur Direktori Baru

```
kijunpos-go/
├── cmd/                    # Command-line entry points
├── config/                 # Konfigurasi aplikasi
├── internal/               # Kode private aplikasi
│   ├── app/                # Application bootstrapping dan DI
│   ├── domain/             # Domain entities dan interfaces
│   ├── usecase/            # Business logic
│   ├── repository/         # Data access implementations
│   ├── delivery/           # Delivery mechanisms (gRPC, HTTP, etc)
│   │   └── grpc/           # gRPC handlers
│   └── pkg/                # Internal shared packages
│       ├── apm/            # Application monitoring
│       └── logger/         # Logging
├── pkg/                    # Public shared packages (jika ada)
└── proto/                  # Proto definitions
```

## Langkah-langkah Migrasi

### 1. Migrasi Bertahap

Untuk saat ini, direktori `app/` belum bisa sepenuhnya dihapus karena:

1. Beberapa komponen sudah dimigrasikan ke struktur baru:
   - `app/lib/apm` → `internal/pkg/apm`
   - `app/lib/logger` → `internal/pkg/logger`

2. Komponen yang masih perlu dimigrasikan:
   - Model-model lain di `app/model`
   - Handler lain di `app/handler`
   - Service lain di `app/service`
   - Repository lain di `app/repository`

### 2. Pendekatan Migrasi yang Disarankan

1. **Migrasi per Domain**:
   - Identifikasi domain bisnis (misalnya: user, product, order)
   - Migrasikan satu domain lengkap dari model hingga handler
   - Uji fungsionalitas domain tersebut
   - Lanjutkan ke domain berikutnya

2. **Langkah-langkah Detail**:
   - Pindahkan entity/model ke `internal/domain/`
   - Definisikan interface repository dan usecase di `internal/domain/`
   - Implementasikan repository di `internal/repository/`
   - Implementasikan usecase di `internal/usecase/`
   - Implementasikan delivery handler di `internal/delivery/grpc/`
   - Update dependency injection di `internal/app/app.go`

### 3. Setelah Migrasi Selesai

Setelah semua komponen berhasil dimigrasikan:

1. Pastikan semua tes berjalan dengan baik
2. Hapus direktori `app/` yang lama
3. Update dokumentasi dan README

## Keuntungan Clean Architecture

1. **Pemisahan Tanggung Jawab**: Setiap layer memiliki tanggung jawab yang jelas
2. **Dependency Rule**: Layer dalam (domain) tidak bergantung pada layer luar
3. **Testability**: Lebih mudah untuk menulis unit test
4. **Maintainability**: Lebih mudah untuk memelihara dan mengembangkan kode
5. **Flexibility**: Lebih mudah untuk mengganti framework atau library eksternal
