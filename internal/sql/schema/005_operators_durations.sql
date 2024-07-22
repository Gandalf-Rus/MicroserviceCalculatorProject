-- +goose Up

CREATE TABLE operators_durations (
    operator_name TEXT PRIMARY KEY NOT NULL,
    operator_duration FLOAT NOT NULL
);

-- +goose Down
DROP TABLE operators_durations;