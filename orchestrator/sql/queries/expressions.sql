-- name: CreateExpression :one
INSERT INTO expressions (id, expression_body, expression_status_id, expression_result)
VALUES($1, $2, $3, $4)
RETURNING *;