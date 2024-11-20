CREATE TABLE redeem_log (
    id TEXT PRIMARY KEY,
    serial TEXT NOT NULL,
    course_id TEXT NOT NULL,
    google_classroom_student_id TEXT NOT NULL,
    assignment_id TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by TEXT NOT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

-- name: GetRedeemLogsByStudentId :many
SELECT 
    id,
    serial,
    course_id,
    google_classroom_student_id,
    assignment_id
FROM redeem_log
WHERE google_classroom_student_id = @google_classroom_student_id;

-- name: CreateRedeemLog :exec
INSERT INTO redeem_log (
    id,
    serial,
    course_id,
    google_classroom_student_id,
    assignment_id,
    created_at,
    created_by
) VALUES (
    @id,
    @serial,
    @course_id,
    @google_classroom_student_id,
    @assignment_id,
    time('now'),
    time('now')
);