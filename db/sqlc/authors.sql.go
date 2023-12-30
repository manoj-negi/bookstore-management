// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: authors.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (
    name,
    is_deleted
) VALUES (
    $1,
    $2
) RETURNING id, name, is_deleted, created_at, updated_at
`

type CreateAuthorParams struct {
	Name      string      `json:"name"`
	IsDeleted pgtype.Bool `json:"is_deleted"`
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error) {
	row := q.db.QueryRow(ctx, createAuthor, arg.Name, arg.IsDeleted)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAuthor = `-- name: DeleteAuthor :one
DELETE FROM authors WHERE id = $1
RETURNING id, name, is_deleted, created_at, updated_at
`

func (q *Queries) DeleteAuthor(ctx context.Context, id int32) (Author, error) {
	row := q.db.QueryRow(ctx, deleteAuthor, id)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllAuthors = `-- name: GetAllAuthors :many
SELECT id, name, is_deleted, created_at, updated_at FROM authors LIMIT 10
`

func (q *Queries) GetAllAuthors(ctx context.Context) ([]Author, error) {
	rows, err := q.db.Query(ctx, getAllAuthors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Author{}
	for rows.Next() {
		var i Author
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IsDeleted,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, name, is_deleted, created_at, updated_at FROM authors WHERE id = $1
`

func (q *Queries) GetAuthor(ctx context.Context, id int32) (Author, error) {
	row := q.db.QueryRow(ctx, getAuthor, id)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateAuthor = `-- name: UpdateAuthor :one
UPDATE authors
SET
    name = CASE 
    WHEN $1::boolean = TRUE THEN $2
    ELSE name
    END,
    is_deleted = CASE
    WHEN $3::boolean = TRUE THEN $4
    ELSE is_deleted
    END
WHERE id = $5
RETURNING id, name, is_deleted, created_at, updated_at
`

type UpdateAuthorParams struct {
	SetName      bool        `json:"set_name"`
	Name         string      `json:"name"`
	SetIsDeleted bool        `json:"set_is_deleted"`
	IsDeleted    pgtype.Bool `json:"is_deleted"`
	ID           int32       `json:"id"`
}

func (q *Queries) UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) (Author, error) {
	row := q.db.QueryRow(ctx, updateAuthor,
		arg.SetName,
		arg.Name,
		arg.SetIsDeleted,
		arg.IsDeleted,
		arg.ID,
	)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
