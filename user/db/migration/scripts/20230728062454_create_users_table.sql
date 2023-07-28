-- +goose Up
CREATE TABLE users (
    id CHAR(16) PRIMARY KEY,
    the_name VARCHAR(256) NOT NULL,
    email VARCHAR(256) NOT NULL UNIQUE,
    the_password VARCHAR(1024) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;
