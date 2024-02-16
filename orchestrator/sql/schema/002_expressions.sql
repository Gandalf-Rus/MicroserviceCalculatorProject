-- +goose Up

CREATE TABLE expressions (
    id TEXT PRIMARY KEY NOT NULL,
    expression_body TEXT UNIQUE NOT NULL,
    expression_status_id INT REFERENCES statuses(id) NOT NULL,
    expression_result FLOAT
);

-- +goose Down
DROP TABLE expressions;