-- name: CreateRefreshTokens :one
insert into refresh_tokens (id, token, user_id, created_at, expires_at)
values ($1, $2, $3, $4, $5) returning *;

-- name: GetRefreshTokenByToken :one
select * from refresh_tokens where token = $1;

-- name: DeleteAllRefreshTokens :exec
delete from refresh_tokens where user_id = $1;

-- name: RefreshTokenExists :one
select exists (select 1 from refresh_tokens where token = $1);
