-- name: CreateOrder :one
INSERT INTO orders (
    book_id,
    user_id,
    order_no,
    quantity,
    total_price,
    status,
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

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: GetAllOrders :many
SELECT * FROM orders LIMIT 10;

-- name: UpdateOrder :one
UPDATE orders
SET
    book_id =  CASE
    WHEN @set_book_id::boolean = TRUE THEN @book_id
    ELSE book_id
    END,
    user_id =  CASE
    WHEN @set_user_id::boolean = TRUE THEN @user_id
    ELSE user_id
    END,
    order_no =  CASE
    WHEN @set_order_no::boolean = TRUE THEN @order_no
    ELSE order_no
    END,
    quantity =  CASE
    WHEN @set_quantity::boolean = TRUE THEN @quantity
    ELSE quantity
    END,
    total_price =  CASE
    WHEN @set_total_price::boolean = TRUE THEN @total_price
    ELSE total_price
    END,
    status =  CASE
    WHEN @set_status::boolean = TRUE THEN @status
    ELSE status
    END,
    is_deleted =  CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeleteOrder :one
DELETE FROM orders WHERE id = $1
RETURNING *;
