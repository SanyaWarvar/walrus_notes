CREATE TABLE IF NOT EXISTS movies(
    id uuid primary key,
    title varchar not null,
    description varchar,
    img_url varchar
);

CREATE TABLE genres(
    id uuid primary key,
    name varchar not null
);

CREATE TABLE IF NOT EXISTS movie_genre(
    movie_id uuid references movies(id),
    genre_id uuid references genres(id)
);