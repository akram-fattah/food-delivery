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



CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    category_name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE meals (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL,
    image_url TEXT,
	is_available BOOLEAN DEFAULT TRUE,
    category_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_category
        FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE SET NULL
);

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    total_price NUMERIC(10,2) NOT NULL DEFAULT 0,
    address TEXT,
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);


CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    meal_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    price NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_meal
        FOREIGN KEY (meal_id)
        REFERENCES meals(id)
        ON DELETE CASCADE
);