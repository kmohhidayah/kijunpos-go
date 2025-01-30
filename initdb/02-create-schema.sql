-- File ini akan otomatis dijalankan terhadap database kijundb
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Masukkan data dummy
INSERT INTO users (username, password_hash, email) 
VALUES
    ('admin', '$2a$10$examplehash', 'admin@example.com'),
    ('user1', '$2a$10$examplehash', 'user1@example.com')
ON CONFLICT (username) DO NOTHING;
