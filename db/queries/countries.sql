-- name: CreateCountry :one
INSERT INTO countries (
    iso2,
    short_name,
    long_name,
    numcode,
    calling_code,
    cctld,
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

-- name: GetCountry :one
SELECT * FROM countries WHERE id = $1;

-- name: GetAllCountries :many
SELECT * FROM countries LIMIT 10;

-- name: UpdateCountry :one
UPDATE countries
SET
    iso2 =CASE
    WHEN @set_iso2::boolean = TRUE THEN @iso2
    ELSE iso2
    END,
    short_name = CASE
    WHEN @set_short_name::boolean = TRUE THEN @short_name
    ELSE short_name
    END,
    long_name = CASE
    WHEN @set_long_name::boolean = TRUE THEN @long_name
    ELSE long_name
    END,
    numcode = CASE
    WHEN @set_numcode::boolean = TRUE THEN @numcode
    ELSE numcode
    END,
    calling_code = CASE
    WHEN @set_calling_code::boolean = TRUE THEN @calling_code
    ELSE calling_code
    END,
    cctld = CASE
    WHEN @set_cctld::boolean = TRUE THEN @cctld
    ELSE cctld
    END,
    is_deleted = CASE
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeleteCountry :one
DELETE FROM countries WHERE id = $1
RETURNING *;
