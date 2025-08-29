CREATE TYPE USER_ROLE AS ENUM ('guest', 'customer', 'employee', 'admin');

CREATE TABLE users (
    id INT PRIMARY KEY,
    username VARCHAR(64) UNIQUE NOT NULL,
    password_hash VARCHAR(128) NOT NULL,
    role USER_ROLE NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted TIMESTAMP
);

CREATE TABLE sessions (
    session_id VARCHAR(64) NOT NULL,
    user_id INT NOT NULL,
    expiry TIMESTAMP,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE components (
    id INT PRIMARY KEY
);