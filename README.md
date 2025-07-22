# Marketplace API

Этот сервис предоставляет REST API для размещения и просмотра объявлений.

## Возможности
1. Регистрация и авторизация пользователей через JWT
2. Создание объявлений с заголовком, описанием, ценой и изображением
3. Получение ленты объявлений с фильтрацией, пагинацией и сортировкой

## Стек технологийл
Язык: Go
БД: PostgreSQL
Контейнеризация: Docker + Docker Compose

## Установка и запуск
1. Клонируй репозиторий:
```bash
git clone https://github.com/Mutter0815/vk-test-task.git
cd vk-test-task
```
2. В файле .env заполните переменные:
```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=marketplace_db

APP_PORT=8080
JWT_SECRET=ваш_длинный_секрет_не_менее_32_символов
```
## Основные команды
Сборка и запуск с помощью Makefile:
```
1. Сборка образов
make build

2. Запуск базы и приложения
make up

3. Применение миграций
make migrate

4. Просмотр логов
make logs

5. Остановка и удаление томов
make down
```
Аналогичные команды без Makefile:
```bash
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d
docker-compose run --rm migrate
docker-compose logs -f app
docker-compose down -v
```

## HTTP API
Базовый адрес: http://localhost:8080

### Регистрация
```bash
POST /auth/register
Content-Type: application/json

{
  "username": "user1",
  "password": "Pass1234"
}

```
Пример curl:
```bash
curl -i -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"Pass1234"}'
```
Успешный ответ:
```bash
{
  "id": 1,
  "username": "user1",
  "message": "Пользователь зарегистрирован"
}
```
### Авторизация
```bash
POST /auth/login
Content-Type: application/json

{
  "username": "user1",
  "password": "Pass1234"
}
```
Пример curl:
```bash
curl -i -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"Pass1234"}'
```
Успешный ответ:
```bash
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
### Создание объявления
```bash
POST /ads
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>

{
  "title":       "Продам велосипед",
  "description": "Городской, в отличном состоянии",
  "image_url":   "https://example.com/bike.jpg",
  "price":       15000
}
```
Пример curl:
```bash
curl -i -X POST http://localhost:8080/ads \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Продам велосипед","description":"Отличное состояние","image_url":"https://example.com/bike.jpg","price":15000}'
```
Успешный ответ:
```bash
{
  "id": 1,
  "user_id": 1,
  "title": "Продам велосипед",
  "description": "Городской, в отличном состоянии",
  "image_url": "https://example.com/bike.jpg",
  "price": 15000,
  "created_at": "2025-07-21T12:34:56Z",
  "isMine": true
}
```
### Просмотр объявлений
```bash
GET /ads
```
Пример curl без авторизации:
```bash
curl -i http://localhost:8080/ads
```
```bash
curl -i http://localhost:8080/ads \
  -H "Authorization: Bearer $TOKEN"
```
Параметры запроса:
- `page` (int, default=1)
- `page_size` (int, default=10)
- `price_min` (uint)
- `price_max` (uint)
- `sort_by` (`date` или `price`)
- `order` (`asc` или `desc`)

Пример:
```bash
curl -i "http://localhost:8080/ads?page=2&page_size=5&sort_by=price&order=asc&price_min=10000"
```
## Работа с базой

Войти в psql внутри контейнера:
```bash
docker exec -it marketplace-db psql -U postgres -d marketplace_db
```
Основные команды psql:
```sql
\dt               -- список таблиц
\d+ users         -- описание таблицы users
SELECT * FROM ads;
```
