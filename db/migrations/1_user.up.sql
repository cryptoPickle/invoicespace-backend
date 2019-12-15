CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users(
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR (5000) NOT NULL,
    last_name VARCHAR (5000) NOT NULL,
    email VARCHAR (3000) NOT NULL UNIQUE,
    password VARCHAR (5000) NOT NULL,
    organisation_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at integer ,
    disabled BOOLEAN ,
    role int
)