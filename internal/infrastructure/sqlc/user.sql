CREATE TABLE users (
    id TEXT PRIMARY KEY,
    course_id TEXT NOT NULL,
    google_classroom_student_id TEXT NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by TEXT NOT NULL,
    updated_at DATETIME,
    updated_by TEXT
);

-- name: CheckUserExistsByEmail :one
SELECT 
    id,
    email,
    created_at,
    created_by,
    updated_at,
    updated_by
FROM users
WHERE email = @email;

-- name: RegisterStudent :exec
INSERT INTO users (
    id,
    course_id,
    google_classroom_student_id,
    firstname,
    lastname,
    email,
    password,
    created_at,
    created_by
) VALUES (
    @id,
    @course_id,
    @google_classroom_student_id,
    @firstname,
    @lastname,
    @email,
    @password,
    time('now'),
    time('now')
);

-- name: GetUserByEmail :one
SELECT 
    id,
    email,
    password,
    created_at,
    created_by,
    updated_at,
    updated_by
FROM users
WHERE email = @email;

-- name: GetUserProfileByEmail :one
SELECT 
    id,
    course_id,
    google_classroom_student_id,
    firstname,
    lastname,
    email,
    created_at,
    created_by,
    updated_at,
    updated_by
FROM users
WHERE email = @email;

-- name: GetUsers :many
SELECT 
    id,
    email,
    created_at,
    created_by,
    updated_at,
    updated_by
FROM users;

-- name: GetUserById :one
SELECT 
    id,
    email,
    created_at,
    created_by,
    updated_at,
    updated_by
FROM users
WHERE id = @id;

-- name: CreateUser :exec
INSERT INTO users (
    id,
    email,
    password,
    created_at,
    created_by
) VALUES (
    @id,
    @email,
    @password,
    NOW(),
    NOW()
);
