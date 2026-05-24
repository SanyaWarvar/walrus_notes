# Walrus Notes API

Backend для заметок с REST API, WebSocket-уведомлениями и совместным редактированием.

- Swagger: `/swagger/index.html`
- Healthcheck: `HEAD /healthz`
- Префикс API: `/wn`

## Локальный запуск

```bash
go run ./cmd/app
```

Конфиг: `config/main.yaml`. Секреты и переопределения — в `./etc/secrets/.env` (или `.env` в корне).

## Деплой

### VPS (Docker Compose)

Postgres, Redis, приложение и nginx поднимаются вместе:

```bash
cp .env.example .env   # опционально
docker compose up -d --build
```

Сервис доступен на порту **80**. Nginx проксирует запросы в Go-приложение и WebSocket.

Остановка:

```bash
docker compose down
```

## WebSocket

### Эндпоинты

| Метод | URL | Назначение |
|-------|-----|------------|
| `GET` | `/wn/api/connection?user_id={uuid}` | Основное соединение: JSON-сообщения, ивенты, черновики |
| `GET` | `/wn/api/room/{noteId}?user_id={uuid}` | Комната совместного редактирования заметки |
| `POST` | `/wn/api/secret` | Генерация секрета для соединения (заготовка) |

`user_id` — UUID пользователя из query-параметра.

### Формат сообщений

Все сообщения — JSON-объект:

```json
{
  "event": "EVENT_NAME",
  "payload": { }
}
```

#### `/wn/api/connection`

- Протокол: **текстовые** WebSocket-фреймы (JSON)
- Сервер каждые **10 секунд** шлёт `PING` — клиент должен ответить:

```json
{ "event": "PONG", "payload": {} }
```

#### `/wn/api/room/{noteId}`

- Протокол: **бинарные** WebSocket-фреймы
- Сообщения проксируются между всеми участниками комнаты (noteId = ID комнаты)
- Текстовые фреймы игнорируются

---

## Ивенты клиента → сервер

Отправляются через `/wn/api/connection`.

### Обновление черновика

> Чтобы отменить черновик — отправь ивент с пустым полем `newDraft`.

Запрос:

```json
{
  "event": "UPDATE_DRAFT_REQUEST",
  "payload": {
    "noteId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
    "newDraft": "текст черновика"
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

### Фиксация черновика

Запрос:

```json
{
  "event": "COMMIT_DRAFT_REQUEST",
  "payload": {
    "noteId": "a78756cf-9d47-4c16-a8d6-17d3207447b4"
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

### Keep-alive

Ответ на серверный `PING`:

```json
{
  "event": "PONG",
  "payload": {}
}
```

---

## Ивенты сервера → клиент

Сервер пушит их на активное соединение `/wn/api/connection`, когда меняются данные (заметки, лейауты, граф).

### Удаление лейаута

```json
{
  "event": "DELETE_LAYOUT",
  "payload": {
    "layoutId": "a78756cf-9d47-4c16-a8d6-17d3207447b4"
  }
}
```

### Обновление лейаута

```json
{
  "event": "UPDATE_LAYOUT",
  "payload": {
    "layoutId": "a78756cf-9d47-4c16-a8d6-17d3207447b4"
  }
}
```

### Создание заметки

```json
{
  "event": "CREATE_NOTE",
  "payload": {
    "layoutId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
    "noteId": "128756cf-9d47-4c16-a8d6-17d3207447b4"
  }
}
```

### Удаление заметки

```json
{
  "event": "DELETE_NOTE",
  "payload": {
    "noteId": "a78756cf-9d47-4c16-a8d6-17d3207447b4"
  }
}
```

### Обновление заметки

```json
{
  "event": "UPDATE_NOTE",
  "payload": {
    "noteId": "a78756cf-9d47-4c16-a8d6-17d3207447b4"
  }
}
```

### Изменение позиции на графе

```json
{
  "event": "CHANGE_NOTE_POSITION",
  "payload": {
    "layoutId": "a78756cf-9d47-4c16-a8d6-17d3207447b4"
  }
}
```

### Изменение связей на графе

```json
{
  "event": "CHANGE_NOTE_LINKS",
  "payload": {
    "layoutId": "a78756cf-9d47-4c16-a8d6-17d3207447b4"
  }
}
```

### Перемещение заметки между лейаутами

```json
{
  "event": "DRAG_NOTE",
  "payload": {
    "noteId": "a78756cf-9d47-4c16-a8d6-17d3207447b4",
    "toLayoutId": "b88756cf-9d47-4c16-a8d6-17d3207447b4"
  }
}
```