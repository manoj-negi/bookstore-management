-- name: CreateCategory :one
INSERT INTO categories (
    name,
    is_special,
    is_deleted
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories WHERE id = $1;

-- name: GetAllCategories :many
SELECT * FROM categories LIMIT 10;

-- name: UpdateCategory :one
UPDATE categories
SET
    name = CASE
    WHEN @set_name::boolean = TRUE THEN @name
    ELSE name
    END,
    is_special = CASE
    WHEN @set_is_special::boolean = TRUE THEN @is_special
    ELSE is_special
    END,
    is_deleted = CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeleteCategory :one
DELETE FROM categories WHERE id = $1
RETURNING *;
