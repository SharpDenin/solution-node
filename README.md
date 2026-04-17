# Solution Node

## Как запустить:

### 1.Устанавливаем git на рабочую станцию (место где вы хотите развернуть проект). Ссылку на скачивание и инструкцию можно найти по ссылке (https://git-scm.com/install/).

### 2. Открываем папку проекта в терминале (Windows + R -> "cmd" -> Enter -> cd /путь к папке)

### 3. Клонируем репозиторий git clone https://github.com/SharpDenin/solution-node.git и ожидаем окончание завершения загрузки

## Альтернативный вариант, разворачивать из архива

### 3.1 Распаковать архив в то место, где хотите развернуть приложение

### После этого можно переходить к шагу №4

### 4. Создаем файл .env в папке проека (там, где лежит docker-compose.yaml)

#### Пример содержания:
.env
```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=report_db


DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=report_db


SERVER_PORT=8080

JWT_SECRET=0387088c2e979e63459319164ad9cc0651d81d63d9fcc662c8cfb6d71099812c // Обязательно замените на свой! (Можно с сайта https://jwtsecrets.com/)

UPLOAD_DIR=./uploads
BASE_URL=http://localhost:8090

CORS_ALLOWED_ORIGINS=http://localhost:3010,http://127.0.0.1:3010
CORS_ALLOW_CREDENTIALS=true
```

### 5. Открываем терминал в корне проекта

### 6. Вводим команду ```docker-compose up --build```

### 7.Если докер не установлен:

a. Установите Docker и Docker Compose:

b. Windows / macOS: Docker Desktop

c. Linux: sudo apt install docker.io docker-compose (или через официальный репозиторий)

d. После установки повторите шаг 3.

## Переход к фронту по ссылке: http://адрес_хоста:3010
## Креды авторизации под админом (Логин: admin, Пароль: snAdmin01)
