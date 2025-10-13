create table if not exists notes(
	id uuid primary key,
	title varchar,
	payload text,
	created_at timestamptz,
	owner_id uuid not null,
	have_access uuid[]
);

create table if not exists layouts(
	id uuid primary key,
	title varchar,
	owner_id uuid not null,
	have_access uuid[]
);

create table if not exists layout_note(
	note_id uuid references notes(id),
	layout_id uuid references layouts(id),
	x_position DOUBLE PRECISION,
	y_position DOUBLE precision,
	primary key (note_id, layout_id)
);

create table if not exists links(
	id uuid primary key,
    layout_id uuid,
	x1_position DOUBLE PRECISION,
	y1_position DOUBLE precision,
	x2_position DOUBLE PRECISION,
	y2_position DOUBLE precision,
	color varchar
);  
