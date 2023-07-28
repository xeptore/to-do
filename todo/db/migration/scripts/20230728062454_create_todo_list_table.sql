-- +goose Up
CREATE TABLE todo_lists (
    id CHAR(16) PRIMARY KEY,
    the_name VARCHAR(256) NOT NULL,
    the_description VARCHAR(1024) NOT NULL,
    created_by_id VARCHAR(16) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE todo_lists;
