-- name: GetDurations :many
SELECT * FROM operators_durations;

-- name: GetDurationByName :one
SELECT * FROM operators_durations
WHERE operator_name = $1;

-- name: EditDuration :one
UPDATE operators_durations
SET operator_duration = $2
WHERE operator_name = $1
RETURNING *;