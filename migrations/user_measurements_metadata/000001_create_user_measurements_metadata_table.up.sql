CREATE TABLE IF NOT EXISTS users_measurements_metadata (
    id int references users(user_id) NOT NULL,
    date date DEFAULT CURRENT_DATE,
    weight float NOT NULL,
    height float NOT NULL,
    neck int,
    shoulders int,
    chest int,
    left_bicep int,
    right_bicep int,
    left_forearm int,
    right_forearm int,
    waist int,
    hips int,
    left_thigh int,
    right_thigh int,
    left_calf int,
    right_calf int
);

ALTER TABLE users_measurements_metadata ADD CONSTRAINT measurements_pk PRIMARY KEY (id, date);