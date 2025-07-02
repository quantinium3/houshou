-- name: CreateEmailVerificationToken :one
insert into email_verification (id, token, user_id, created_at, expires_at)
    values ($1, $2, $3, $4, $5) returning *;

-- name: GetEmailVerificationByToken :one
select * from email_verification where token = $1;

-- name: GetEmailVerificationTokenById :one
select * from email_verification where user_id = $1;

-- name: DeleteAllEmailVerificationToken :exec
delete from email_verification where user_id = $1;