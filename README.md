# Social Network Otus

## Локальный запуск

1. Убедитесь, что установлены Docker, Docker Compose, git, postman.
2. Склонируйте репозиторий:
   ```bash
   git clone https://github.com/Stalles/Otus.git
   cd socialNetworkOtus
   ```
3. При необходимости (скорее всего, не нужно) поменяйте енвы бд в docker-compose.yml (см. в корне проекта).
4. Запустите сервисы:
   ```bash
   docker-compose up --build -d
   ```
5. Приложение будет доступно на `http://localhost:8080`.

## Тестирование API

Для тестирования используйте Postman-коллекцию `postman_collection.json` (см. в корне проекта).

### Эндпоинты:
- `POST /user/register` — регистрация пользователя
- `POST /login` — авторизация
- `GET /user/get/{id}` — получение анкеты по ID

## Postman-коллекция
См. файл `postman_collection.json` в корне проекта.