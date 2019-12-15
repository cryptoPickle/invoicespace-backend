CREATE TABLE token_revoke_list(
                      token_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                      refreshToken VARCHAR (1000) NOT NULL,
                      user_id UUID REFERENCES users(user_id) NOT NULL
)