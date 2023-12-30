-- name: CreateRole :one
INSERT INTO roles (
    name,
    description,
    is_deleted
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetRole :one
SELECT * FROM roles WHERE id = $1;

-- name: GetAllRoles :many
SELECT * FROM roles LIMIT 10;

-- name: UpdateRole :one
UPDATE roles
SET
    name = CASE
    WHEN @set_name::boolean = TRUE THEN @name
    ELSE name
    END,
    description = CASE
    WHEN @set_description::boolean = TRUE THEN @description
    ELSE description
    END,
    is_deleted = CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeleteRole :one
DELETE FROM roles WHERE id = $1
RETURNING *;
