// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: expressions.sql

package database

import (
	"context"
	"database/sql"
)

const createExpression = `-- name: CreateExpression :one
INSERT INTO expressions (id, expression_body, expression_status_id, count_of_subexpression, expression_result)
VALUES($1, $2, 2, $3, NULL)
RETURNING id, expression_body, expression_status_id, count_of_subexpression, expression_result
`

type CreateExpressionParams struct {
	ID                   string
	ExpressionBody       string
	CountOfSubexpression int32
}

func (q *Queries) CreateExpression(ctx context.Context, arg CreateExpressionParams) (Expression, error) {
	row := q.db.QueryRowContext(ctx, createExpression, arg.ID, arg.ExpressionBody, arg.CountOfSubexpression)
	var i Expression
	err := row.Scan(
		&i.ID,
		&i.ExpressionBody,
		&i.ExpressionStatusID,
		&i.CountOfSubexpression,
		&i.ExpressionResult,
	)
	return i, err
}

const editExpressions = `-- name: EditExpressions :one
UPDATE expressions
SET expression_status_id = $3,
    expression_result = $2
WHERE id = $1
RETURNING id, expression_body, expression_status_id, count_of_subexpression, expression_result
`

type EditExpressionsParams struct {
	ID                 string
	ExpressionResult   sql.NullFloat64
	ExpressionStatusID int32
}

func (q *Queries) EditExpressions(ctx context.Context, arg EditExpressionsParams) (Expression, error) {
	row := q.db.QueryRowContext(ctx, editExpressions, arg.ID, arg.ExpressionResult, arg.ExpressionStatusID)
	var i Expression
	err := row.Scan(
		&i.ID,
		&i.ExpressionBody,
		&i.ExpressionStatusID,
		&i.CountOfSubexpression,
		&i.ExpressionResult,
	)
	return i, err
}

const getExpressionByID = `-- name: GetExpressionByID :one
SELECT id, expression_body, expression_status_id, count_of_subexpression, expression_result FROM expressions WHERE id = $1
`

func (q *Queries) GetExpressionByID(ctx context.Context, id string) (Expression, error) {
	row := q.db.QueryRowContext(ctx, getExpressionByID, id)
	var i Expression
	err := row.Scan(
		&i.ID,
		&i.ExpressionBody,
		&i.ExpressionStatusID,
		&i.CountOfSubexpression,
		&i.ExpressionResult,
	)
	return i, err
}
