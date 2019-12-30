CREATE TABLE user_pool(
    id UUID REFERENCES user_pools(id),
    userId varchar(50) NOT NULL,
    user_role INT NOT NULL
)