-- +goose up
create table users (
    id text primary key,
    name text not null unique,
    email text not null unique,
    password text not null,
    subdomain text not null,
    email_verified boolean not null default false,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table refresh_tokens (
    id text primary key,
    token text unique not null,
    user_id text not null,
    created_at timestamp not null default current_timestamp,
    expires_at timestamp not null,
    foreign key (user_id) references users(id) on delete cascade on update restrict
);

create table reset_tokens (
    id text primary key,
    token text unique not null,
    user_id text not null,
    created_at timestamp not null default current_timestamp,
    expires_at timestamp not null,
    foreign key (user_id) references users(id) on delete cascade on update restrict
);

create table email_verification (
    id text primary key,
    token text unique not null,
    user_id text not null,
    created_at timestamp not null default current_timestamp,
    expires_at timestamp not null,
    foreign key (user_id) references users(id) on delete cascade on update restrict
);

-- +goose down
drop table refresh_tokens;
drop table reset_tokens;
drop table email_verification;
drop table users;