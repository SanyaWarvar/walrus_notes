# Сокет ивенты

## Обновление черновика
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
        "status"
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
