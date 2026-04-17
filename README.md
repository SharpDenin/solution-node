# solution-node

# Как запустить:
## 1. Создаем файл .env в папке проека (там, где лежит docker-compose.yaml)
### Пример содержания:
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

JWT_SECRET=0387088c2e909e63439319064ad8cc0641d81d63d9fcc662c8cfb6d71099812c

UPLOAD_DIR=./uploads
BASE_URL=http://localhost:8090

CORS_ALLOWED_ORIGINS=http://localhost:3010,http://127.0.0.1:3010
CORS_ALLOW_CREDENTIALS=true
```

## 2. Открываем терминал в корне проекта
## 3. Вводим команду ```docker-compose up --build```

