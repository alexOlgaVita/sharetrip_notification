# job4j_notification

### О проекте

Проект реализует сервис "Notification", который отвечает за хранение и отправку уведомлений.
В дальнейшем будет использоваться в качестве микросервсиного приложения в проекте "Совместные поездки" (наряду с другими имкросервисами - такими, как "ShareTrip"").

### Стек технологий

GoLang 1.25.0, Fiber, PostgreSQL

### Используемые инструменты

fmt, linter, тесты

### Требования к окружению

GoLang 1.25.0, PostgreSQL

#### API
Первая версия сервиса реализует REST API contract:

|      |                     |                             |
|------|---------------------|-----------------------------|
| POST | /notifications      | Создание уведомления        |
| GET  | /notifications/{id} | Получение уведомления       |


Пример ответа POST-запроса "/notifications":
{
"recipient_id": "client-123",
"type": "trip_published",
"payload": {
"trip_id": "trip-456"
}

Пример ответа GET-запроса на получение увемоления "/notifications/{id}":
{
"id": "notification-789",
"recipient_id": "client-123",
"type": "trip_published",
"status": "created",
"payload": {
"trip_id": "trip-456"
},
"created_at": "2026-07-06T12:00:00Z"
}

### Контакты

![Ильина Ольга](images/olga.jpg)

- Telegram: [@OlgaIlyina0312](https://t.me/OlgaIlyina0312)
- Email:    [oliljina@mail.ru](oliljina@mail.ru)