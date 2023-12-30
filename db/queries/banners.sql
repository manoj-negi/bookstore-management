-- name: CreateBanner :one
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
) RETURNING *;

-- name: GetBanner :one
SELECT * FROM banners WHERE id = $1;

-- name: GetAllBanners :many
SELECT * FROM banners LIMIT 10;

-- name: UpdateBanner :one
UPDATE banners
SET
    name = CASE
    WHEN @set_name::boolean = TRUE THEN @name
    ELSE name
    END,
    image = CASE
    WHEN @set_image::boolean = TRUE THEN @image
    ELSE image
    END,
    start_date = CASE
    WHEN @set_start_date::boolean = TRUE THEN @start_date
    ELSE start_date
    END,
    end_date = CASE
    WHEN @set_end_date::boolean = TRUE THEN @end_date
    ELSE end_date
    END,
    offer_id = CASE
    WHEN @set_offer_id::boolean = TRUE THEN @offer_id
    ELSE offer_id
    END,
    is_deleted = CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeleteBanner :one
DELETE FROM banners WHERE id = $1
RETURNING *;
