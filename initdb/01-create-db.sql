DO
$do$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'kijun_user_db') THEN
      CREATE USER kijun_user_db WITH PASSWORD 'admin123';
   END IF;
END
$do$;

-- Buat database jika belum ada
DO
$do$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'kijundb') THEN
      CREATE DATABASE kijundb OWNER kijun_user_db;
   END IF;
END
$do$;

-- Berikan semua hak akses kepada user
GRANT ALL PRIVILEGES ON DATABASE kijundb TO kijun_user_db;
