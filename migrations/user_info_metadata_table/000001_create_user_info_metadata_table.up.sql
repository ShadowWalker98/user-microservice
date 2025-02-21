CREATE TABLE users_info_metadata (
    id int references users(user_id) NOT NULL PRIMARY KEY,
    phone text,
    dob date,
    profile text
);