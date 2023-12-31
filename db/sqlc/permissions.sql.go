// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: permissions.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPermission = `-- name: CreatePermission :one
INSERT INTO permissions (
    name,
    permission,
    is_deleted
) VALUES (
    $1,
    $2,
    $3
) RETURNING id, name, permission, is_deleted, created_at, updated_at
`

type CreatePermissionParams struct {
	Name       string      `json:"name"`
	Permission pgtype.Text `json:"permission"`
	IsDeleted  pgtype.Bool `json:"is_deleted"`
}

func (q *Queries) CreatePermission(ctx context.Context, arg CreatePermissionParams) (Permission, error) {
	row := q.db.QueryRow(ctx, createPermission, arg.Name, arg.Permission, arg.IsDeleted)
	var i Permission
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Permission,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePermission = `-- name: DeletePermission :one
DELETE FROM permissions WHERE id = $1
RETURNING id, name, permission, is_deleted, created_at, updated_at
`

func (q *Queries) DeletePermission(ctx context.Context, id int32) (Permission, error) {
	row := q.db.QueryRow(ctx, deletePermission, id)
	var i Permission
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Permission,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllPermissions = `-- name: GetAllPermissions :many
SELECT id, name, permission, is_deleted, created_at, updated_at FROM permissions LIMIT 10
`

func (q *Queries) GetAllPermissions(ctx context.Context) ([]Permission, error) {
	rows, err := q.db.Query(ctx, getAllPermissions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Permission{}
	for rows.Next() {
		var i Permission
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Permission,
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

const getPermission = `-- name: GetPermission :one
SELECT id, name, permission, is_deleted, created_at, updated_at FROM permissions WHERE id = $1
`

func (q *Queries) GetPermission(ctx context.Context, id int32) (Permission, error) {
	row := q.db.QueryRow(ctx, getPermission, id)
	var i Permission
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Permission,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePermission = `-- name: UpdatePermission :one
UPDATE permissions
SET
    name = CASE
    WHEN $1::boolean = TRUE THEN $2
    ELSE name
    END,
    permission = CASE
    WHEN $3::boolean = TRUE THEN $4
    ELSE permission
    END,
    is_deleted = CASE
    WHEN $5::boolean = TRUE THEN $6
    ELSE is_deleted
    END
WHERE id = $7
RETURNING id, name, permission, is_deleted, created_at, updated_at
`

type UpdatePermissionParams struct {
	SetName       bool        `json:"set_name"`
	Name          string      `json:"name"`
	SetPermission bool        `json:"set_permission"`
	Permission    pgtype.Text `json:"permission"`
	SetIsDeleted  bool        `json:"set_is_deleted"`
	IsDeleted     pgtype.Bool `json:"is_deleted"`
	ID            int32       `json:"id"`
}

func (q *Queries) UpdatePermission(ctx context.Context, arg UpdatePermissionParams) (Permission, error) {
	row := q.db.QueryRow(ctx, updatePermission,
		arg.SetName,
		arg.Name,
		arg.SetPermission,
		arg.Permission,
		arg.SetIsDeleted,
		arg.IsDeleted,
		arg.ID,
	)
	var i Permission
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Permission,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
