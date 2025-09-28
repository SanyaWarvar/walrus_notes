DROP TABLE IF EXISTS roles;
CREATE TABLE IF NOT EXISTS roles(
    id int primary key,
    name varchar not null,
    permissions jsonb
);

INSERT INTO roles VALUES
(1, 'ADMINT'),
(2, 'CLIENT');

CREATE TABLE IF NOT EXISTS users(
    id uuid primary key,
    username varchar unique not null,
    email varchar unique not null,
    password varchar not null,
    role varchar not null,
    img_url varchar not null,
    confirmed_email boolean default 'f',
    created_at timestamptz not null
);