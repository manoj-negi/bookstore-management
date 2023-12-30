-- name: CreateOffer :one
INSERT INTO offers (
    book_id,
    discount_percentage,
    start_date,
    end_date,
    is_deleted
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: GetOffer :one
SELECT * FROM offers WHERE id = $1;

-- name: GetAllOffers :many
SELECT * FROM offers LIMIT 10;

-- name: UpdateOffer :one
UPDATE offers
SET
    book_id = CASE
    WHEN @set_book_id::boolean = TRUE THEN @book_id
    ELSE book_id
    END,
    discount_percentage = CASE
    WHEN @set_discount_percentage::boolean = TRUE THEN @discount_percentage
    ELSE discount_percentage
    END,
    start_date = CASE
    WHEN @set_start_date::boolean = TRUE THEN @start_date
    ELSE start_date
    END,
    end_date = CASE
    WHEN @set_end_date::boolean = TRUE THEN @end_date
    ELSE end_date
    END,
    is_deleted = CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeleteOffer :one
DELETE FROM offers WHERE id = $1
RETURNING *;
