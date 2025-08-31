CREATE TYPE USER_ROLE AS ENUM ('guest', 'customer', 'employee', 'admin');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(32) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role USER_ROLE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE sessions (
    session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT UNIQUE NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);


-- add default user
-- [user:pass]
-- admin:admin
-- test:test
INSERT INTO "users" ("id", "email", "username", "password_hash", "role", "created_at", "deleted_at") VALUES
(1,	'admin@gmail.com',	'admin',	'$2a$10$9RL0UWVE29y57k4QohuId.2KxVfhxJ4.nwihulpdn4JPIpheFmMgC',	'admin',	'2025-08-31 10:20:00.511205',	NULL),
(2,	'test@gmail.com',	'test',	'$2a$10$fbvP1s2Je02iOkpinVo1ZO6ishUFMji9DyvQYI2T.T5cGHOLflSge',	'guest',	'2025-08-31 10:20:39.121888',	NULL);
