drop table links;

create table if not exists links(
    layout_id uuid references layouts(id),
    first_note_id uuid references notes(id),
    second_note_id uuid references notes(id),
    primary key(layout_id, first_note_id, second_note_id)
);