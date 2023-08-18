CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    root_directory VARCHAR(255) NOT NULL
);