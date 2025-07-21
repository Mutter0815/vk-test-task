# Marketplace API

Этот сервис предоставляет REST API для размещения и просмотра объявлений.

## Запуск

```bash
make build
make up
```

После запуска приложение доступно на `http://localhost:8080`.
Swagger-документация доступна по адресу `http://localhost:8080/swagger/index.html`.

## Аутентификация

Для выполнения защищённых запросов необходимо передавать JWT токен в заголовке `Authorization` вида:

```
Authorization: Bearer <token>
```

Получить токен можно через эндпоинт `/auth/login`.
