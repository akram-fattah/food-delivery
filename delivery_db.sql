-- First Create postgres database named "delivery_db" and then run this SQL script to create the necessary tables.

DROP TRIGGER IF EXISTS set_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column;
DROP TABLE IF EXISTS users;

 CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,

    username VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,

    email VARCHAR(150) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,

    phone VARCHAR(20) UNIQUE,
    address TEXT,

    role VARCHAR(20) NOT NULL DEFAULT 'user',

    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verification_code VARCHAR(100),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_role CHECK (role IN ('user', 'admin', 'Delivery'))
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


ALTER TABLE users
ADD COLUMN verification_expires TIMESTAMP;