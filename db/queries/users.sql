-- name: CreateUser :one
INSERT INTO users (
    first_name,
    last_name,
    gender,
    dob,
    address,
    city,
    state,
    country_id,
    mobile_no,
    username,
    email,
    password,
    role_id,
    otp,
    is_deleted
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetAllUsers :many
SELECT * FROM users LIMIT 10;

-- name: UpdateUser :one
UPDATE users
SET
    first_name = CASE 
    WHEN @set_first_name::boolean = TRUE THEN @first_name
    ELSE first_name
    END,
    last_name = CASE 
    WHEN @set_last_name::boolean = TRUE THEN @last_name
    ELSE last_name
    END,
    gender = CASE 
    WHEN @set_gender::boolean = TRUE THEN @gender
    ELSE gender
    END,
    dob = CASE 
    WHEN @set_dob::boolean = TRUE THEN @dob
    ELSE dob
    END,
    address = CASE 
    WHEN @set_address::boolean = TRUE THEN @address
    ELSE address
    END,
    city = CASE 
    WHEN @set_city::boolean = TRUE THEN @city
    ELSE city
    END,
    state = CASE 
    WHEN @set_state::boolean = TRUE THEN @state
    ELSE state
    END,
    country_id = CASE 
    WHEN @set_country_id::boolean = TRUE THEN @country_id
    ELSE country_id
    END,
    mobile_no = CASE 
    WHEN @set_mobile_no::boolean = TRUE THEN @mobile_no
    ELSE mobile_no
    END,
    username = CASE 
    WHEN @set_username::boolean = TRUE THEN @username
    ELSE username
    END,
    email = CASE 
    WHEN @set_email::boolean = TRUE THEN @email
    ELSE email
    END,
    password = CASE 
    WHEN @set_password::boolean = TRUE THEN @password
    ELSE password
    END,
    role_id = CASE 
    WHEN @set_role_id::boolean = TRUE THEN @role_id
    ELSE role_id
    END,
    otp = CASE 
    WHEN @set_otp::boolean = TRUE THEN @otp
    ELSE otp
    END,
    is_deleted = CASE 
    WHEN @set_is_deleted::boolean = TRUE THEN @is_deleted
    ELSE is_deleted
    END
WHERE id = @id
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users WHERE id = $1
RETURNING *;
