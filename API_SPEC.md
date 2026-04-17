# Описание API

## Общая информация

- Базовый URL: `http://localhost:8090` (пример, замените на актуальный)
- Формат данных: JSON (кроме загрузки файлов и экспорта)
- Авторизация: для защищённых эндпоинтов (`🔐`) требуется JWT-токен, полученный при входе. Токен передаётся в заголовке:
    - `Authorization: Bearer <JWT_TOKEN>`
- Роли пользователей:
    - `worker` – может создавать отчёты, загружать файлы.
    - `admin` – имеет доступ ко всем отчётам, управляет вопросами, экспортирует данные.

---

## 🔐 Аутентификация

### Регистрация пользователя

Создаёт нового пользователя с ролью `worker`. Пароль хешируется с помощью bcrypt.

**POST** `/register`

**Тело запроса:**
```json
{
  "full_name": "John Doe",
  "login": "john",
  "password": "123456"
}
```

**Ответ:**
```json
201 Created
```

### Вход в систему

**POST** `/login`

**Тело запроса:**
```json
{
  "login": "john",
  "password": "123456"
}
```

**Ответ:**
```json
200 OK
```
```json
{
  "token": "JWT_TOKEN"
}
```

> В JWT-токене содержатся `user_id` и `role`.

### Информация о текущем пользователе

**GET** `/me` 🔐

**Заголовки:** `Authorization: Bearer JWT_TOKEN`

**Ответ:**
```json
200 OK
```
Текст: `Hello user {userID} role: {role}`

---

## 📤 Загрузка файлов

### Загрузить файл

**POST** `/upload` 🔐  
`Content-Type: multipart/form-data`

**Тело запроса:** поле `file` (бинарные данные)

**Ответ:**
```json
200 OK
```
```json
{
  "url": "http://localhost:8090/uploads/filename.jpg"
}
```

---

## 📋 Отчёты

### Создать отчёт (только worker)

**POST** `/reports` 🔐

**Тело запроса:**
```json
{
  "place": "Plant 1",
  "report_date": "2026-04-16",
  "responsible_name": "John Doe",
  "answers": [
    {
      "question_id": "uuid",
      "text": "All good",
      "image_url": "http://..."
    }
  ]
}
```

**Ответ:**
```json
201 Created
```

> `user_id` берётся из JWT-токена.

### Получить список отчётов (только admin)

**GET** `/reports` 🔐

**Параметры запроса (query):**
- `date_from` – начальная дата (YYYY-MM-DD)
- `date_to` – конечная дата (YYYY-MM-DD)
- `place` – фильтр по месту
- `user_id` – UUID пользователя
- `limit` – количество записей
- `offset` – смещение для пагинации

Пример: `/reports?date_from=2026-01-01&date_to=2026-12-31&limit=10&offset=0`

**Ответ:**
```json
200 OK
```
```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "place": "Plant 1",
    "report_date": "2026-04-16",
    "responsible_name": "John",
    "created_at": "2026-04-16T10:00:00Z"
  }
]
```

### Получить детали отчёта (только admin)

**GET** `/reports/{id}` 🔐

**Параметры пути:** `id` – UUID отчёта.

**Ответ:**
```json
200 OK
```
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "place": "Plant 1",
  "report_date": "2026-04-16",
  "responsible_name": "John",
  "created_at": "2026-04-16T10:00:00Z",
  "answers": [
    {
      "question_id": "uuid",
      "question_text": "Is everything ok?",
      "answer_text": "Yes",
      "image_url": "http://..."
    }
  ]
}
```

### Экспорт отчётов в Excel (только admin)

**GET** `/reports/export` 🔐

**Параметры запроса (необязательные):** `date_from`, `date_to`, `place`, `user_id`

**Ответ:**
```json
200 OK
```
Файл Excel (`.xlsx`)  
Content-Type: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`

---

## ❓ Управление вопросами (только admin)

### Получить список вопросов

**GET** `/questions` 🔐

**Ответ:**
```json
200 OK
```
```json
[
  {
    "id": "uuid",
    "text": "Is everything ok?",
    "order_index": 1,
    "is_active": true
  }
]
```

### Создать вопрос

**POST** `/questions` 🔐

**Тело запроса:**
```json
{
  "text": "Is everything ok?",
  "order_index": 1
}
```

**Ответ:**
```json
201 Created
```

> По умолчанию `is_active = true`.

### Обновить вопрос

**PUT** `/questions/{id}` 🔐

**Параметры пути:** `id` – UUID вопроса.

**Тело запроса:**
```json
{
  "text": "Updated question",
  "order_index": 2,
  "is_active": true
}
```

**Ответ:**
```json
200 OK
```

### Удалить вопрос (soft delete)

**DELETE** `/questions/{id}` 🔐

**Параметры пути:** `id` – UUID вопроса.

**Ответ:**
```json
200 OK
```

> Вопрос не удаляется физически, а помечается как неактивный (`is_active = false`).