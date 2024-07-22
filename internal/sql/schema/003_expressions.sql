-- +goose Up

CREATE TABLE expressions (
    id TEXT PRIMARY KEY NOT NULL,
    expression_body TEXT UNIQUE NOT NULL,
    expression_status_id INT REFERENCES statuses(id) NOT NULL DEFAULT(2),
    count_of_subexpression INT NOT NULL,
    expression_result FLOAT DEFAULT(NULL)
);

-- +goose Down
DROP TABLE expressions;