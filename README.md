# Сокет ивенты

### Обновление черновика
> **(!)** Чтобы отменить создание черновика - необходимо отправить ивент с пустым полем "newDraft"

Отправить:
```json
{
    "event": "UPDATE_DRAFT_REQUEST",
    "payload": {
        "noteId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
        "newDraft": "11"
    }
}
```
Ответ:
```json
{
    "event": "UPDATE_DRAFT_RESPONSE",
    "payload": {
        "status": "true"
    }
}
```
***

```json
{
    "event": "COMMIT_DRAFT_REQUEST",
    "payload": {
        "noteId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
    }
}
```
Ответ:
```json
{
    "event": "COMMIT_DRAFT_RESPONSE",
    "payload": {
        "status": "true"
    }
}
```
***

## Ивенты сервера

### Удаление лейаута

```json
{
    "event": "DELETE_LAYOUT",
    "payload": {
        "layoutId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
    }
}
```

### Обновление лейаута

```json
{
    "event": "UPDATE_LAYOUT",
    "payload": {
        "layoutId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
    }
}
```

### Создание заметки

```json
{
    "event": "CREATE_NOTE",
    "payload": {
        "layoutId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
        "noteId": "128756cf-9d47-4c16-a8d6-17d3207447b4",
    }
}
```

### Удаление заметки

```json
{
    "event": "DELETE_NOTE",
    "payload": {
        "noteId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
    }
}
```

### Обновление заметки

```json
{
    "event": "UPDATE_NOTE",
    "payload": {
        "noteId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
    }
}
```

### Изменение графа

```json
{
    "event": "CHANGE_NOTE_POSITION",
    "payload": {
        "noteId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
    }
}
```

### Изменение связи

```json
{
    "event": "CHANGE_NOTE_LINKS",
    "payload": {
        "noteId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
    }
}
```

### Перемещение заметки между папок

```json
{
    "event": "DRAG_NOTE",
    "payload": {
        "noteId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
    }
}
```