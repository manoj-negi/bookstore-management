-- name: CreateAuthor :one
INSERT INTO authors (
    name,
    is_deleted
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: GetAuthor :one
SELECT * FROM authors WHERE id = $1;

-- name: GetAllAuthors :many
SELECT * FROM authors LIMIT 10;

-- name: UpdateAuthor :one
UPDATE authors
SET
    name = CASE 
    WHEN @set_name::boolean = TRUE THEN @name
    ELSE name
    END,
    is_deleted = CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeleteAuthor :one
DELETE FROM authors WHERE id = $1
RETURNING *;
