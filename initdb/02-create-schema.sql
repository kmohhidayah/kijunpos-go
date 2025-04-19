-- File ini akan otomatis dijalankan terhadap database kijundb
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE,
    whatsapp_number VARCHAR(20) UNIQUE,
    otp_pin VARCHAR(6),
    is_active BOOLEAN NOT NULL DEFAULT true,
    failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,
    password_changed_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT chk_contact_info CHECK (email IS NOT NULL OR whatsapp_number IS NOT NULL)
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_whatsapp ON users(whatsapp_number);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Hapus data yang mungkin sudah ada untuk menghindari konflik
TRUNCATE TABLE users;

-- Dummy data untuk login dengan email/password
-- Password: password123 (bcrypt hash)
INSERT INTO users (
    id,
    username, 
    password_hash, 
    email, 
    is_active, 
    failed_login_attempts,
    created_at,
    password_changed_at
) 
VALUES
    (
        '11111111-1111-1111-1111-111111111111',
        'email_user',
        '$2a$10$lT.Lx2GsvtRYEfVpwGfh8e9HM8MJW8.eLm6Ar.iBhstGBxUclfAPO',
        'email_user@example.com',
        true,
        0,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '22222222-2222-2222-2222-222222222222',
        'admin',
        '$2a$10$lT.Lx2GsvtRYEfVpwGfh8e9HM8MJW8.eLm6Ar.iBhstGBxUclfAPO',
        'admin@example.com',
        true,
        0,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

-- Dummy data untuk login dengan WhatsApp/PIN
INSERT INTO users (
    id,
    username, 
    password_hash,
    whatsapp_number,
    pin,
    is_active, 
    failed_login_attempts,
    created_at
) 
VALUES
    (
        '33333333-3333-3333-3333-333333333333',
        'whatsapp_user',
        '',
        '+6281234567890',
        '123456',
        true,
        0,
        CURRENT_TIMESTAMP
    ),
    (
        '44444444-4444-4444-4444-444444444444',
        'dual_auth_user',
        '$2a$10$lT.Lx2GsvtRYEfVpwGfh8e9HM8MJW8.eLm6Ar.iBhstGBxUclfAPO',
        '+6289876543210',
        '654321',
        true,
        0,
        CURRENT_TIMESTAMP
    );

-- User dengan email dan WhatsApp (bisa login dengan kedua metode)
INSERT INTO users (
    id,
    username, 
    password_hash,
    email,
    whatsapp_number,
    pin,
    is_active, 
    failed_login_attempts,
    created_at
) 
VALUES
    (
        '55555555-5555-5555-5555-555555555555',
        'complete_user',
        '$2a$10$lT.Lx2GsvtRYEfVpwGfh8e9HM8MJW8.eLm6Ar.iBhstGBxUclfAPO',
        'complete_user@example.com',
        '+6287654321098',
        '111222',
        true,
        0,
        CURRENT_TIMESTAMP
    );
