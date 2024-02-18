-- name: CreateExpression :one
INSERT INTO expressions (id, expression_body, expression_status_id, count_of_subexpression, expression_result)
VALUES($1, $2, 2, $3, NULL)
RETURNING *;

-- name: GetExpressionByID :one
SELECT * FROM expressions WHERE id = $1;