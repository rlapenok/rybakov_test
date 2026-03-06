# Withdrawal Service

REST API сервис для создания заявок на вывод средств с поддержкой идемпотентности и защитой от двойного списания.

---

## Быстрый старт

### Требования

- Docker + Docker Compose

### Запуск

```bash
# 1. Скопировать файл с переменными окружения
cp .env.example .env

# 2. Поднять сервис (БД + миграции + приложение)
docker compose up --build -d
```

Приложение будет доступно на `http://localhost:8080`.  
Миграции применяются автоматически при старте.

### Остановка

```bash
docker compose down
```

---

## Переменные окружения (`.env`)

| Переменная | Описание | По умолчанию |
|---|---|---|
| `HTTP_SERVER_PORT` | Порт HTTP-сервера | `8080` |
| `AUTH_BEARER_TOKEN` | Статический Bearer-токен для авторизации | `dev-token` |
| `LOGGER_LEVEL` | Уровень логирования (`debug`, `info`, `warn`, `error`) | `debug` |
| `LOGGER_FORMAT` | Формат логов (`console`, `json`) | `console` |
| `PG_HOST` | Хост PostgreSQL | `postgres` |
| `PG_PORT` | Порт PostgreSQL | `5432` |
| `PG_USER` | Пользователь БД | `postgres` |
| `PG_PASSWORD` | Пароль БД | `postgres` |
| `PG_DATABASE` | Имя БД | `rybakov` |
| `PG_MIGRATION_PATH` | Путь к папке с миграциями | `migrations` |

---

## API

Все эндпоинты требуют заголовок:
```
Authorization: Bearer <AUTH_BEARER_TOKEN>
```

### `POST /api/v1/withdrawals`

Создать заявку на вывод средств.

**Request:**
```json
{
  "user_id":        "123e4567-e89b-12d3-a456-426614174000",
  "amount":         "50.00",
  "currency":       "USDT",
  "destination":    "123e4567-e89b-12d3-a456-426614174001",
  "idempotency_key": "unique-key-per-request"
}
```

**Responses:**

| Код | Причина |
|---|---|
| `201` | Withdrawal создан, возвращает `{"id": "uuid"}` |
| `400` | Невалидные данные (amount ≤ 0, неверная валюта и т.д.) |
| `401` | Отсутствует или неверный Bearer-токен |
| `404` | Пользователь не найден |
| `409` | Недостаточно средств на балансе |
| `422` | Тот же `idempotency_key` с другим payload |

**Идемпотентность:**  
Повторный запрос с тем же `idempotency_key` и тем же payload вернёт `201` с тем же ID без повторного списания.

---

### `GET /api/v1/withdrawals/:id`

Получить информацию о заявке по ID.

**Response `200`:**
```json
{
  "id":             "uuid",
  "user_id":        "uuid",
  "amount":         "50.00",
  "currency":       "USDT",
  "destination":    "uuid",
  "idempotency_key": "unique-key-per-request",
  "status":         "pending"
}
```

| Код | Причина |
|---|---|
| `200` | Данные withdrawal |
| `401` | Неверный токен |
| `404` | Withdrawal не найден |

---

## Запуск тестов

Тесты интеграционные, запускают PostgreSQL в Docker через testcontainers.  
Требуется работающий Docker daemon.

```bash
go test ./internal/infra/repository/... -v -timeout 120s
```

**Что покрывают тесты:**

| Тест | Сценарий |
|---|---|
| `TestCreateWithdrawal_Success` | Успешное создание withdrawal |
| `TestCreateWithdrawal_InsufficientBalance` | Баланс меньше суммы → `409` |
| `TestCreateWithdrawal_Idempotency` | Один ключ + одинаковый payload → тот же ID; одинаковый ключ + разный payload → `422` |
| `TestCreateWithdrawal_Concurrent` | Два параллельных запроса на один баланс — ровно одно списание |

---

## Ключевые решения

### Защита от двойного списания (консистентность)

Создание withdrawal выполняется в одной транзакции PostgreSQL:

1. **`SELECT balance FROM users WHERE id = $1 FOR UPDATE`** — пессимистическая блокировка строки пользователя. Все конкурентные запросы к тому же `user_id` ожидают завершения текущей транзакции.
2. `INSERT INTO withdrawals ... ON CONFLICT (user_id, idempotency_key) DO NOTHING` — вставка с защитой от дублей.
3. `UPDATE users SET balance = balance - $amount WHERE id = $user_id AND balance >= $amount` — атомарное списание только при наличии достаточного баланса.
4. `COMMIT` — либо оба действия применяются, либо откатываются вместе.

Дополнительная страховка на уровне БД: `CHECK (balance >= 0)` — не позволит записать отрицательный баланс даже при ошибке в логике приложения.

### Идемпотентность

- Constraint `UNIQUE (user_id, idempotency_key)` в БД — физически не позволяет создать дубль.
- `payload_hash` (SHA-256 от `user_id|amount|destination`) сохраняется вместе с withdrawal.
- При повторном запросе с тем же ключом:
  - hash совпал → возвращается исходный ID без повторного списания.
  - hash отличается → `422 IDEMPOTENCY_PAYLOAD_MISMATCH`.
- Область ключа — пользователь (`user_id`), а не глобальная, чтобы ключи разных пользователей не конфликтовали.

### Безопасность

- Фиксированный Bearer-токен из env — базовая auth без утечки данных сессий.
- Внутренние ошибки (`500`) не возвращают stack trace или детали реализации в ответе.
- Валидация входных данных на уровне доменных типов до обращения к БД.
