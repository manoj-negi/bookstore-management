-- name: CreateCategoryImage :one
INSERT INTO categories_images (
    image,
    category_id,
    is_deleted
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetCategoryImage :one
SELECT * FROM categories_images WHERE id = $1;

-- name: GetAllCategoryImages :many
SELECT * FROM categories_images LIMIT 10;

-- name: UpdateCategoryImage :one
UPDATE categories_images
SET
image =  CASE
    WHEN @set_image::boolean = TRUE THEN @image
    ELSE image
    END,
    category_id =  CASE
    WHEN @set_category_id::boolean = TRUE THEN @category_id
    ELSE category_id
    END,
    is_deleted =  CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeleteCategoryImage :one
DELETE FROM categories_images WHERE id = $1
RETURNING *;
