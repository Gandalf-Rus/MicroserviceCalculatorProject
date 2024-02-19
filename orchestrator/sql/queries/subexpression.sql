-- name: CreateSubexpression :one
INSERT INTO subexpressions (
    expression_id, 
    subexpression_number, 
    subexpression_body, 
    subexpression_status_id, 
    subexpression_result)
VALUES($1, $2, $3, 2, NULL)
RETURNING *;


-- name: GetSubexpressionByExprID :many
SELECT * FROM subexpressions WHERE expression_id = $1;

-- name: GetSubexpressionByStatusID :many
SELECT * FROM subexpressions WHERE subexpression_status_id = $1;

-- name: GetSubexpressionByExprIDAndNumber :one
SELECT * FROM subexpressions WHERE expression_id = $1 AND subexpression_number = $2;

-- name: EditSubexpressions :one
UPDATE subexpressions
SET subexpression_status_id = $3,
    subexpression_result = $4
WHERE expression_id = $1 AND subexpression_number = $2
RETURNING *;

-- name: EditSubexpressionStatus :one
UPDATE subexpressions
SET subexpression_status_id = $3
WHERE expression_id = $1 AND subexpression_number = $2
RETURNING *;