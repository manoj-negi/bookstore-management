// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: banners.sql

package db

import (
	"context"
	"time"
	"github.com/jackc/pgx/v5/pgtype"
)

const createBanner = `-- name: CreateBanner :one
INSERT INTO banners (
    name,
    image,
    start_date,
    end_date,
    offer_id,
    is_deleted
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING id, name, image, start_date, end_date, offer_id, is_deleted, created_at, updated_at
`

type CreateBannerParams struct {
	Name      string      `json:"name"`
	Image     string 	  `json:"image"`
	StartDate time.Time   `json:"start_date"`
	EndDate   time.Time   `json:"end_date"`
	OfferID   int32       `json:"offer_id"`
	IsDeleted pgtype.Bool `json:"is_deleted"`
}

func (q *Queries) CreateBanner(ctx context.Context, arg CreateBannerParams) (Banner, error) {
	row := q.db.QueryRow(ctx, createBanner,
		arg.Name,
		arg.Image,
		arg.StartDate,
		arg.EndDate,
		arg.OfferID,
		arg.IsDeleted,
	)
	var i Banner
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.StartDate,
		&i.EndDate,
		&i.OfferID,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteBanner = `-- name: DeleteBanner :one
DELETE FROM banners WHERE id = $1
RETURNING id, name, image, start_date, end_date, offer_id, is_deleted, created_at, updated_at
`

func (q *Queries) DeleteBanner(ctx context.Context, id int32) (Banner, error) {
	row := q.db.QueryRow(ctx, deleteBanner, id)
	var i Banner
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.StartDate,
		&i.EndDate,
		&i.OfferID,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllBanners = `-- name: GetAllBanners :many
SELECT id, name, image, start_date, end_date, offer_id, is_deleted, created_at, updated_at FROM banners LIMIT 10
`

func (q *Queries) GetAllBanners(ctx context.Context) ([]Banner, error) {
	rows, err := q.db.Query(ctx, getAllBanners)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Banner{}
	for rows.Next() {
		var i Banner
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Image,
			&i.StartDate,
			&i.EndDate,
			&i.OfferID,
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

const getBanner = `-- name: GetBanner :one
SELECT id, name, image, start_date, end_date, offer_id, is_deleted, created_at, updated_at FROM banners WHERE id = $1
`

func (q *Queries) GetBanner(ctx context.Context, id int32) (Banner, error) {
	row := q.db.QueryRow(ctx, getBanner, id)
	var i Banner
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.StartDate,
		&i.EndDate,
		&i.OfferID,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateBanner = `-- name: UpdateBanner :one
UPDATE banners
SET
    name = CASE
    WHEN $1::boolean = TRUE THEN $2
    ELSE name
    END,
    image = CASE
    WHEN $3::boolean = TRUE THEN $4
    ELSE image
    END,
    start_date = CASE
    WHEN $5::boolean = TRUE THEN $6
    ELSE start_date
    END,
    end_date = CASE
    WHEN $7::boolean = TRUE THEN $8
    ELSE end_date
    END,
    offer_id = CASE
    WHEN $9::boolean = TRUE THEN $10
    ELSE offer_id
    END,
    is_deleted = CASE
    WHEN $11::boolean = TRUE THEN $12
    ELSE is_deleted
    END
WHERE id = $13
RETURNING id, name, image, start_date, end_date, offer_id, is_deleted, created_at, updated_at
`

type UpdateBannerParams struct {
	SetName      bool        `json:"set_name"`
	Name         string      `json:"name"`
	SetImage     bool        `json:"set_image"`
	Image        string      `json:"image"`
	SetStartDate bool        `json:"set_start_date"`
	StartDate    time.Time   `json:"start_date"`
	SetEndDate   bool        `json:"set_end_date"`
	EndDate      time.Time   `json:"end_date"`
	SetOfferID   bool        `json:"set_offer_id"`
	OfferID      int32       `json:"offer_id"`
	SetIsDeleted bool        `json:"set_is_deleted"`
	IsDeleted    pgtype.Bool `json:"is_deleted"`
	ID           int32       `json:"id"`
}

func (q *Queries) UpdateBanner(ctx context.Context, arg UpdateBannerParams) (Banner, error) {
	row := q.db.QueryRow(ctx, updateBanner,
		arg.SetName,
		arg.Name,
		arg.SetImage,
		arg.Image,
		arg.SetStartDate,
		arg.StartDate,
		arg.SetEndDate,
		arg.EndDate,
		arg.SetOfferID,
		arg.OfferID,
		arg.SetIsDeleted,
		arg.IsDeleted,
		arg.ID,
	)
	var i Banner
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.StartDate,
		&i.EndDate,
		&i.OfferID,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
