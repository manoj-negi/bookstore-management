-- name: CreateBook :one
INSERT INTO books (
    title,
    author_id,
    publication_date,
    price,
    stock_quantity,
    bestseller,
    is_deleted
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
) RETURNING *;

-- name: GetBook :one
SELECT * FROM books WHERE id = $1;

-- name: GetAllBooks :many
SELECT * FROM books LIMIT 10;

-- name: GetBestseller :many
SELECT * FROM books WHERE bestseller = true;

-- name: UpdateBook :one
UPDATE books
SET
    title = CASE
    WHEN @set_title::boolean = TRUE THEN @title
    ELSE title
    END,
    author_id = CASE
    WHEN @set_author_id::boolean = TRUE THEN @author_id
    ELSE author_id
    END,
    publication_date = CASE
    WHEN @set_publication_date::boolean = TRUE THEN @publication_date
    ELSE publication_date
    END,
    price = CASE
    WHEN @set_price::boolean = TRUE THEN @price
    ELSE price
    END,
    stock_quantity = CASE
    WHEN @set_stock_quantity::boolean = TRUE THEN @stock_quantity
    ELSE stock_quantity
    END,
    bestseller = CASE
    WHEN @set_bestseller::boolean = TRUE THEN @bestseller
    ELSE bestseller
    END,
    is_deleted = CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeleteBook :one
DELETE FROM books WHERE id = $1
RETURNING *;
