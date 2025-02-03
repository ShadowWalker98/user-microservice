CREATE TABLE users_water_serving_data (
    id int references users(user_id),
    water_serving_size int NOT NULL,
    date date DEFAULT CURRENT_DATE,
    required_servings int NOT NULL,
    consumed int NOT NULL
);

ALTER TABLE users_water_serving_data ADD CONSTRAINT users_water_pk PRIMARY KEY (id, date);