CREATE TABLE users (
    user_id BIGSERIAL PRIMARY KEY,
    email text NOT NULL UNIQUE,
    username text NOT NULL UNIQUE,
    password text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

