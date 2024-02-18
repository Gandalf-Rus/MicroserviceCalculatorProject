-- +goose Up
INSERT INTO statuses (status_name) VALUES 
    ('done'),
    ('in_progress'),
    ('uncorrect');

-- +goose Down
DELETE FROM statuses WHERE status_name IN ('done', 'in_progress', 'uncorrect');
