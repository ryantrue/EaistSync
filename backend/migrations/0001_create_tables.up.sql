CREATE TABLE IF NOT EXISTS contracts (
                                         id BIGINT PRIMARY KEY,
                                         data JSONB
);

CREATE TABLE IF NOT EXISTS states (
                                      id BIGINT PRIMARY KEY,
                                      data JSONB
);
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );