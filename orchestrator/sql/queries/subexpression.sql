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