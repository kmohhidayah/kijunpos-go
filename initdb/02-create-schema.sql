-- File ini akan otomatis dijalankan terhadap database kijundb
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,
    password_changed_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

INSERT INTO users (
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
        'admin',
        '$2a$10$examplehash',
        'admin@example.com',
        true,
        0,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        'user1',
        '$2a$10$examplehash',
        'user1@example.com',
        true,
        0,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    )
ON CONFLICT (username) DO NOTHING;
