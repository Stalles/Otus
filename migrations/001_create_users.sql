CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    second_name VARCHAR(100) NOT NULL,
    birth_date DATE NOT NULL,
    biography TEXT,
    city VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL
); 