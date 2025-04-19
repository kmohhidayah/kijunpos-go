-- Catatan: Sebagian besar inisialisasi database sudah ditangani oleh variabel lingkungan Docker Compose:
-- POSTGRES_USER, POSTGRES_PASSWORD, dan POSTGRES_DB
-- Script ini hanya untuk memberikan hak akses tambahan dan pengaturan lainnya

-- Berikan semua hak akses kepada user
GRANT ALL PRIVILEGES ON DATABASE kijundb TO kijun_user_db;

-- Aktifkan ekstensi yang diperlukan
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
