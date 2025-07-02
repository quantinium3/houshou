-- name: CreateUser :one
insert into users (id, name, email, password, subdomain, email_verified)
values ($1, $2, $3, $4, $5, $6) returning *;

-- name: GetUserById :one
select * from users where id = $1 limit 1;

-- name: UserExists :one
select exists (select 1 from users where email = $1);

-- name: UpdateUserPassword :exec
update users set password = $1 where id = $2;

-- name: UpdateEmailVerificationStatus :exec
update users set email_verified = $1 where id = $2;

-- name: GetUserByEmail :one
select * from users where email = $1 limit 1;