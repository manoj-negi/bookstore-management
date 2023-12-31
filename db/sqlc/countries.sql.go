// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: countries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createCountry = `-- name: CreateCountry :one
INSERT INTO countries (
    iso2,
    short_name,
    long_name,
    numcode,
    calling_code,
    cctld,
    is_deleted
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
) RETURNING id, iso2, short_name, long_name, numcode, calling_code, cctld, is_deleted, created_at, updated_at
`

type CreateCountryParams struct {
	Iso2        string      `json:"iso2"`
	ShortName   string      `json:"short_name"`
	LongName    string      `json:"long_name"`
	Numcode     pgtype.Text `json:"numcode"`
	CallingCode string      `json:"calling_code"`
	Cctld       string      `json:"cctld"`
	IsDeleted   pgtype.Bool `json:"is_deleted"`
}

func (q *Queries) CreateCountry(ctx context.Context, arg CreateCountryParams) (Country, error) {
	row := q.db.QueryRow(ctx, createCountry,
		arg.Iso2,
		arg.ShortName,
		arg.LongName,
		arg.Numcode,
		arg.CallingCode,
		arg.Cctld,
		arg.IsDeleted,
	)
	var i Country
	err := row.Scan(
		&i.ID,
		&i.Iso2,
		&i.ShortName,
		&i.LongName,
		&i.Numcode,
		&i.CallingCode,
		&i.Cctld,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteCountry = `-- name: DeleteCountry :one
DELETE FROM countries WHERE id = $1
RETURNING id, iso2, short_name, long_name, numcode, calling_code, cctld, is_deleted, created_at, updated_at
`

func (q *Queries) DeleteCountry(ctx context.Context, id int32) (Country, error) {
	row := q.db.QueryRow(ctx, deleteCountry, id)
	var i Country
	err := row.Scan(
		&i.ID,
		&i.Iso2,
		&i.ShortName,
		&i.LongName,
		&i.Numcode,
		&i.CallingCode,
		&i.Cctld,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllCountries = `-- name: GetAllCountries :many
SELECT id, iso2, short_name, long_name, numcode, calling_code, cctld, is_deleted, created_at, updated_at FROM countries LIMIT 10
`

func (q *Queries) GetAllCountries(ctx context.Context) ([]Country, error) {
	rows, err := q.db.Query(ctx, getAllCountries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Country{}
	for rows.Next() {
		var i Country
		if err := rows.Scan(
			&i.ID,
			&i.Iso2,
			&i.ShortName,
			&i.LongName,
			&i.Numcode,
			&i.CallingCode,
			&i.Cctld,
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

const getCountry = `-- name: GetCountry :one
SELECT id, iso2, short_name, long_name, numcode, calling_code, cctld, is_deleted, created_at, updated_at FROM countries WHERE id = $1
`

func (q *Queries) GetCountry(ctx context.Context, id int32) (Country, error) {
	row := q.db.QueryRow(ctx, getCountry, id)
	var i Country
	err := row.Scan(
		&i.ID,
		&i.Iso2,
		&i.ShortName,
		&i.LongName,
		&i.Numcode,
		&i.CallingCode,
		&i.Cctld,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateCountry = `-- name: UpdateCountry :one
UPDATE countries
SET
    iso2 =CASE
    WHEN $1::boolean = TRUE THEN $2
    ELSE iso2
    END,
    short_name = CASE
    WHEN $3::boolean = TRUE THEN $4
    ELSE short_name
    END,
    long_name = CASE
    WHEN $5::boolean = TRUE THEN $6
    ELSE long_name
    END,
    numcode = CASE
    WHEN $7::boolean = TRUE THEN $8
    ELSE numcode
    END,
    calling_code = CASE
    WHEN $9::boolean = TRUE THEN $10
    ELSE calling_code
    END,
    cctld = CASE
    WHEN $11::boolean = TRUE THEN $12
    ELSE cctld
    END,
    is_deleted = CASE
    WHEN $13::boolean = TRUE THEN $14
    ELSE is_deleted
    END
WHERE id = $15
RETURNING id, iso2, short_name, long_name, numcode, calling_code, cctld, is_deleted, created_at, updated_at
`

type UpdateCountryParams struct {
	SetIso2        bool        `json:"set_iso2"`
	Iso2           string      `json:"iso2"`
	SetShortName   bool        `json:"set_short_name"`
	ShortName      string      `json:"short_name"`
	SetLongName    bool        `json:"set_long_name"`
	LongName       string      `json:"long_name"`
	SetNumcode     bool        `json:"set_numcode"`
	Numcode        pgtype.Text `json:"numcode"`
	SetCallingCode bool        `json:"set_calling_code"`
	CallingCode    string      `json:"calling_code"`
	SetCctld       bool        `json:"set_cctld"`
	Cctld          string      `json:"cctld"`
	SetIsDeleted   bool        `json:"set_is_deleted"`
	IsDeleted      pgtype.Bool `json:"is_deleted"`
	ID             int32       `json:"id"`
}

func (q *Queries) UpdateCountry(ctx context.Context, arg UpdateCountryParams) (Country, error) {
	row := q.db.QueryRow(ctx, updateCountry,
		arg.SetIso2,
		arg.Iso2,
		arg.SetShortName,
		arg.ShortName,
		arg.SetLongName,
		arg.LongName,
		arg.SetNumcode,
		arg.Numcode,
		arg.SetCallingCode,
		arg.CallingCode,
		arg.SetCctld,
		arg.Cctld,
		arg.SetIsDeleted,
		arg.IsDeleted,
		arg.ID,
	)
	var i Country
	err := row.Scan(
		&i.ID,
		&i.Iso2,
		&i.ShortName,
		&i.LongName,
		&i.Numcode,
		&i.CallingCode,
		&i.Cctld,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
