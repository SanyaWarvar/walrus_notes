create table if not exists permissions (
    id uuid primary key,
    to_user_id uuid,
    from_user_id uuid,
    target_id uuid,
    kind varchar,
    can_read boolean,
    can_write boolean,
    can_edit boolean,
    created_at timestamptz
);