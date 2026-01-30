# Subscription Service

REST API сервис для управления пользовательскими подписками. Реализован на Go с использованием принципов **Clean Architecture**.

Сервис позволяет создавать подписки, обновлять их статус, удалять, получать списки и рассчитывать итоговую стоимость трат пользователя за определенный период.

## Технологический стек

* **Язык:** Go (Golang) 1.23+
* **Архитектура:** Clean Architecture (Handlers -> UseCase -> Storage)
* **Роутер:** [chi](https://github.com/go-chi/chi)
* **База данных:** PostgreSQL
* **Миграции:** [goose](https://github.com/pressly/goose) (embedded)
* **Логирование:** slog (структурированное логирование)
* **Документация:** Swagger (swaggo)
* **Контейнеризация:** Docker, Docker Compose
* **Конфигурация:** YAML + переменные окружения (.env)

## Функционал

* **CRUD подписок:** Создание, чтение, обновление, удаление.
* **Агрегация:** Расчет суммарной стоимости подписок за указанный период (с учетом дат начала и окончания).
* **Docker:** Полная изоляция окружения, запуск одной командой.
* **Swagger UI:** Интерактивная документация API.
* **Graceful Shutdown:** Корректное завершение работы сервера и соединений с БД.

## Быстрый старт (Docker)

Самый простой способ запустить проект — использовать Docker Compose.

1.  **Клонируйте репозиторий:**
    ```bash
    git clone [https://github.com/ZakharMarinin/Subscription-Service.git](https://github.com/ZakharMarinin/Subscription-Service.git)
    cd Subscription-Service
    ```

2.  **Создайте файл `.env`:**
    ```bash
    cp .env.example .env
    # Или создайте вручную со следующим содержимым:
    CONFIG_PATH="./config.yaml"
    POSTGRES_URL={полный адрес подключения}
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=subscription_service
    ```

3.  **Запустите проект:**
    ```bash
    docker-compose up --build
    ```

После запуска сервис будет доступен по адресу: `http://0.0.0.0:8085`

## Документация API (Swagger)

После запуска сервиса документация доступна по адресу:

**[http://0.0.0.0:8085/swagger/index.html](http://0.0.0.0:8085/swagger/index.html)**

Вы можете тестировать запросы прямо из браузера.

### Основные эндпоинты:

* `POST /api/v1/subscriptions` — Создать подписку.
* `GET /api/v1/subscriptions` — Получить список (можно фильтровать по `user_id`).
* `GET /api/v1/subscriptions/total` — Получить сумму трат за период.
* `PUT /api/v1/subscriptions/{id}` — Обновить подписку.
* `DELETE /api/v1/subscriptions/{id}` — Удалить подписку.

## Структура проекта

Проект следует стандарту **Golang Project Layout**:

```text
.
├── cmd/
│   └── main.go             # Точка входа в приложение
├── docs/                   # Сгенерированная документация Swagger
├── internal/
│   ├── application/        # Сборка приложения (App struct)
│   ├── config/             # Логика загрузки конфига
│   ├── domain/             # Основные сущности (Models)
│   ├── http/
│   │   ├── handlers/       # HTTP хендлеры (Transport layer)
│   │   ├── middleware/     # Логгер API запросов
│   │   └── router/         # Настройка маршрутов и middleware
│   ├── storage/            # Работа с базой данных (Repository layer)
│   │   └── migrations/     # SQL файлы миграций
│   └── usecase/            # Бизнес-логика
├── docker-compose.yml
├── Dockerfile
└── README.md