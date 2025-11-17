CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50) NOT NULL,
    enableTOTP BOOLEAN DEFAULT FALSE
);


INSERT INTO users (username, password, email, role)
VALUES 
    ('giang', '123456', 'giang@example.com', 'admin'),
    ('tuan', 'abcdef', 'tuan@example.com', 'user'),
    ('haixon', 'password', 'haixon@example.com', 'user');
