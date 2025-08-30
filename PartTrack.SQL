CREATE TYPE USER_ROLE AS ENUM ('guest', 'customer', 'employee', 'admin');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(128) UNIQUE NOT NULL,
    username VARCHAR(32) NOT NULL,
    password_hash VARCHAR(128) NOT NULL,
    role USER_ROLE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE sessions (
    session_id VARCHAR(64) NOT NULL,
    user_id INT UNIQUE NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE components (
    id SERIAL PRIMARY KEY
);