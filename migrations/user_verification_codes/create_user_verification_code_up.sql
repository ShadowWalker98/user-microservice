CREATE TABLE users_verifications (
    user_id int references users(user_id) PRIMARY KEY,
    verification_code int not null,
    expiry timestamp DEFAULT NOW() + INTERVAL '10 minutes' NOT NULL
)