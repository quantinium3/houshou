// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: users.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
insert into users (id, name, email, password, subdomain, email_verified)
values ($1, $2, $3, $4, $5, $6) returning id, name, email, password, subdomain, email_verified, created_at, updated_at
`

type CreateUserParams struct {
	ID            string
	Name          string
	Email         string
	Password      string
	Subdomain     string
	EmailVerified bool
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.Subdomain,
		arg.EmailVerified,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Subdomain,
		&i.EmailVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
select id, name, email, password, subdomain, email_verified, created_at, updated_at from users where email = $1 limit 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Subdomain,
		&i.EmailVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
select id, name, email, password, subdomain, email_verified, created_at, updated_at from users where id = $1 limit 1
`

func (q *Queries) GetUserById(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Subdomain,
		&i.EmailVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateEmailVerificationStatus = `-- name: UpdateEmailVerificationStatus :exec
update users set email_verified = $1 where id = $2
`

type UpdateEmailVerificationStatusParams struct {
	EmailVerified bool
	ID            string
}

func (q *Queries) UpdateEmailVerificationStatus(ctx context.Context, arg UpdateEmailVerificationStatusParams) error {
	_, err := q.db.Exec(ctx, updateEmailVerificationStatus, arg.EmailVerified, arg.ID)
	return err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
update users set password = $1 where id = $2
`

type UpdateUserPasswordParams struct {
	Password string
	ID       string
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.Password, arg.ID)
	return err
}

const userExists = `-- name: UserExists :one
select exists (select 1 from users where email = $1)
`

func (q *Queries) UserExists(ctx context.Context, email string) (bool, error) {
	row := q.db.QueryRow(ctx, userExists, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
