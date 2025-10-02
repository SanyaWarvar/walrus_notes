create table if not exists notes(
	id uuid primary key,
	title varchar,
	payload text,
	created_at timestamptz,
	owner_id uuid not null,
	have_access uuid[]
);

create table if not exists layout(
	layout_id uuid primary key,
	owner_id uuid not null,
	have_access uuid[]
);

create table if not exists note_position(
	note_id uuid,
	layout_id uuid,
	x_position DOUBLE PRECISION,
	y_position DOUBLE precision,
	primary key (note_id, layout_id)
);

create table if not exists links(
	id uuid primary key,
    layout_id uuid,
	x_position DOUBLE PRECISION,
	y_position DOUBLE precision,
	color varchar
);  
