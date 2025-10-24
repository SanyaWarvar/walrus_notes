drop table links;

create table if not exists links(
    layout_id uuid,
    first_note_id uuid,
    second_note_id uuid,
    primary key(layout_id, first_note_id, second_note_id)
);