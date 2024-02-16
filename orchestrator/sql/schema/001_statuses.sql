-- +goose Up
CREATE TABLE statuses (
    id SERIAL PRIMARY KEY,
    status_name VARCHAR(20) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE statuses;