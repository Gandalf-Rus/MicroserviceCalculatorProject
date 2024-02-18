-- +goose Up

CREATE TABLE subexpressions (
    expression_id TEXT NOT NULL,
    subexpression_number INT NOT NULL,
    subexpression_body TEXT NOT NULL,
    subexpression_status_id INT REFERENCES statuses(id) NOT NULL DEFAULT(2),
    subexpression_result FLOAT DEFAULT(NULL),
    PRIMARY KEY (expression_id, subexpression_number),
    FOREIGN KEY (expression_id) REFERENCES expressions(id)
);

-- +goose Down
DROP TABLE subexpressions;