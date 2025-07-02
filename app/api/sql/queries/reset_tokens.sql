-- name: CreateResetToken :one
insert into reset_tokens (id, token, user_id, created_at, expires_at)
values ($1, $2, $3, $4, $5) returning *;

-- name: GetResetTokenByToken :one
select * from reset_tokens where token = $1;

-- name: DeleteAllResetTokens :exec
delete from reset_tokens where user_id = $1;