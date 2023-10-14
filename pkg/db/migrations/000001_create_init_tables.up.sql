-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

SET TIMEZONE="Europe/Istanbul";

-- Create users table
CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
    email VARCHAR (255) NOT NULL UNIQUE,
    name_surname VARCHAR (100) NOT NULL,
    surname  VARCHAR (50) NOT NULL,
    password_hash VARCHAR (255) NOT NULL,
    user_status INT NOT NULL,
    user_role VARCHAR (25) NOT NULL
);

-- Add indexes
CREATE INDEX active_users ON users (id) WHERE user_status = 1;
