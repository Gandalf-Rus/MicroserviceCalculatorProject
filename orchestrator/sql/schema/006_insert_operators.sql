-- +goose Up
INSERT INTO operators_durations (operator_name, operator_duration) VALUES 
    ('+', 1),
    ('-', 1),
    ('*', 1),
    ('/', 1);

-- +goose Down
DELETE FROM statuses WHERE operator_name IN ('+', '-', '*', '/');
