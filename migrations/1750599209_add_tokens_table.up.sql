CREATE TABLE IF NOT EXISTS refresh_tokens(
    id uuid primary key,
    user_id uuid not null references users(id),
    access_id uuid not null,
    exp_at timestamptz not null
);