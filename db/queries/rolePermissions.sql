-- name: CreateRolePermission :one
INSERT INTO roles_permissions (
    role_id,
    permission_id,
    is_deleted
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetRolePermission :one
SELECT * FROM roles_permissions WHERE id = $1;

-- name: GetAllRolePermissions :many
SELECT * FROM roles_permissions LIMIT 10;

-- name: UpdateRolePermission :one
UPDATE roles_permissions
SET
    role_id = CASE
    WHEN @set_role_id::boolean = TRUE THEN @role_id
    ELSE role_id
    END,
    permission_id = CASE
    WHEN @set_permission_id::boolean = TRUE THEN @permission_id
    ELSE permission_id
    END,
    is_deleted = CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeleteRolePermission :one
DELETE FROM roles_permissions WHERE id = $1
RETURNING *;
