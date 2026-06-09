# Harbor Management System

Веб-приложение для управления морским портом. Позволяет вести учёт судов, отслеживать визиты и управлять доступом сотрудников.

## Стек

- **Backend:** Go, Gin, PostgreSQL
- **Frontend:** React, Vite

## Запуск

### База данных

```bash
psql -U postgres -c "CREATE DATABASE harbor;"
psql -U postgres -d harbor -f Source/backend/schema.sql
```

### Backend

```bash
cd Source/backend
go mod tidy
go run main.go
```

### Frontend

```bash
cd Source/frontend
npm install
npm run dev
```

Фронтенд: `http://localhost:5173`  
Бэкенд: `http://localhost:8080`

## Переменные окружения

Создай `.env` в `Source/backend/`:

```env
USE_DB=true
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=harbor
SERVER_PORT=8080
JWT_SECRET=secret
JWT_EXPIRY_HOURS=24
```

При `USE_DB=false` приложение запустится без базы данных — данные хранятся в памяти.

## Тестовые аккаунты

| Email | Пароль | Роль |
|-------|--------|------|
| admin@harbor.ru | admin123 | Администратор |
| dispatcher@harbor.ru | disp123 | Диспетчер |
| operator@harbor.ru | oper123 | Оператор |

## Структура проекта

```
Source/
├── backend/
│   ├── config/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   ├── service/
│   ├── schema.sql
│   └── main.go
└── frontend/
    └── src/
        ├── api/
        ├── components/
        ├── context/
        └── pages/
```