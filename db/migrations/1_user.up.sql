CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users(
    user_id UUID PRIMARY KEY UNIQUE,
    first_name VARCHAR (50) NOT NULL,
    last_name VARCHAR (50) NOT NULL,
    email VARCHAR (300) NOT NULL,
    password VARCHAR (50) NOT NULL,
    organisation_id UUID,
    created_at integer NOT NULL,
    updated_at integer NOT NULL,
    disabled BOOLEAN NOT NULL,
    role int NOT NULL
)