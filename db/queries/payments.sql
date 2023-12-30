-- name: CreatePayment :one
INSERT INTO payments (
    order_id,
    amount,
    payment_status,
    is_deleted
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments WHERE id = $1;

-- name: GetAllPayments :many
SELECT * FROM payments LIMIT 10;

-- name: UpdatePayment :one
UPDATE payments
SET
    order_id = CASE
    WHEN @set_order_id::boolean = TRUE THEN @order_id
    ELSE order_id
    END,
    amount = CASE
    WHEN @set_amount::boolean = TRUE THEN @amount
    ELSE amount
    END,
    payment_status = CASE
    WHEN @set_payment_status::boolean = TRUE THEN @payment_status
    ELSE payment_status
    END,
    is_deleted = CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeletePayment :one
DELETE FROM payments WHERE id = $1
RETURNING *;
