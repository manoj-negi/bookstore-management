-- name: CreateBookCategory :one
INSERT INTO books_categories (
    book_id,
    category_id,
    is_deleted
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetBookCategory :one
SELECT * FROM books_categories WHERE id = $1;

-- name: GetAllBookCategories :many
SELECT * FROM books_categories LIMIT 10;

-- name: UpdateBookCategory :one
UPDATE books_categories
SET
    book_id = CASE
    WHEN @set_book_id::boolean = TRUE THEN @book_id
    ELSE book_id
    END,
    category_id = CASE
    WHEN @set_category_id::boolean = TRUE THEN @category_id
    ELSE category_id
    END,
    is_deleted = CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeleteBookCategory :one
DELETE FROM books_categories WHERE id = $1
RETURNING *;
