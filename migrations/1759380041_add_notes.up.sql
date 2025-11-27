create table if not exists notes(
	id uuid primary key,
	title varchar,
	payload text,
	created_at timestamptz,
	owner_id uuid not null references users(id),
	have_access uuid[],
	layout_id uuid references layouts(id)
);

create table if not exists layouts(
	id uuid primary key,
	title varchar,
	owner_id uuid not null references users(id),
	have_access uuid[]
);

create table if not exists positions(
	note_id uuid  primary key references notes(id),
	x_position DOUBLE PRECISION,
	y_position DOUBLE PRECISION
);
