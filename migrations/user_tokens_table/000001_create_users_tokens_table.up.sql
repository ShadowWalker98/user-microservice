CREATE TABLE users_tokens (
    user_id int references users(user_id),
    token_type int,
    token text,
    expiry timestamp DEFAULT NOW() + INTERVAL '2 days'
);