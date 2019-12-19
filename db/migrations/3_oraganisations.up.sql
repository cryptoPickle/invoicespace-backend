CREATE TABLE organisations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(50) UNIQUE NOT NULL,
    description varchar(500),
    user_pool_id UUID UNIQUE DEFAULT uuid_generate_v4(),
    worker_limit INT NOT NULL,
    user_limit INT NOT NULL,
    disabled BOOLEAN,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at integer
)