-- name: CreatePermission :one
INSERT INTO permissions (
    name,
    permission,
    is_deleted
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetPermission :one
SELECT * FROM permissions WHERE id = $1;

-- name: GetAllPermissions :many
SELECT * FROM permissions LIMIT 10;

-- name: UpdatePermission :one
UPDATE permissions
SET
    name = CASE
    WHEN @set_name::boolean = TRUE THEN @name
    ELSE name
    END,
    permission = CASE
    WHEN @set_permission::boolean = TRUE THEN @permission
    ELSE permission
    END,
    is_deleted = CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeletePermission :one
DELETE FROM permissions WHERE id = $1
RETURNING *;
