# Сокет ивенты

## Обновление черновика
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
