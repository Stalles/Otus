# Social Network Otus

## Локальный запуск

### Вариант 1: Docker Compose (рекомендуется)

1. Убедитесь, что установлены Docker и Docker Compose.
2. Склонируйте репозиторий:
   ```bash
   git clone <repo_url>
   cd socialNetworkOtus
   ```
3. Создайте файл `.env` в корне проекта (пример ниже).
4. Запустите сервисы:
   ```bash
   docker-compose up --build
   ```
5. Приложение будет доступно на `http://localhost:8080`.

### Вариант 2: Ручной запуск

1. Установите Go 1.21+ и PostgreSQL.
2. Склонируйте репозиторий и перейдите в папку проекта.
3. Создайте файл `.env` с параметрами подключения к БД (пример ниже).
4. Примените миграции:
   ```bash
   task migrate
   ```
5. Запустите приложение:
   ```bash
   task run
   ```

## Пример .env
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=social_network
JWT_SECRET=your_jwt_secret
```

## Тестирование API

Для тестирования используйте Postman-коллекцию `postman_collection.json` (см. в корне проекта).

### Основные эндпоинты:
- `POST /user/register` — регистрация пользователя
- `POST /login` — авторизация
- `GET /user/get/{id}` — получение анкеты по ID

## Postman-коллекция
См. файл `postman_collection.json` в корне проекта.

---

## Taskfile

Для удобства работы с миграциями и запуском приложения используется [Taskfile](https://taskfile.dev/). Установите task, если он ещё не установлен:

```bash
brew install go-task/tap/go-task # для MacOS
# или
sudo snap install go-task --classic # для Linux
```

### Основные команды
- `task migrate` — применить миграции к базе данных
- `task run` — запустить приложение (с авто-миграцией и автозапуском БД) 