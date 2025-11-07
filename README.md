# Personal Finance Ledger API

Это учебный проект: простой API-сервис на Go для учета личных финансов. Сервис позволяет создавать транзакции и устанавливать бюджеты по категориям.

## Как запустить

1.  Клонируйте репозиторий.
2.  Выполните в корневой папке проекта:
    ```bash
    go run ./gateway/cmd/gateway
    ```
3.  Сервер будет запущен и доступен по адресу `http://localhost:8080`.

## Примеры использования API

### Бюджеты

*   `POST /api/budgets` — Создать или обновить бюджет.
*   `GET /api/budgets` — Получить список бюджетов.

**Пример создания бюджета:**
```bash
curl -X POST http://localhost:8080/api/budgets \
-H "Content-Type: application/json" \
-d '{"category":"еда","limit":5000}'
```

### Транзакции

*   `POST /api/transactions` — Создать транзакцию.
*   `GET /api/transactions` — Получить список транзакций.

**Пример создания транзакции:**
```bash
curl -X POST http://localhost:8080/api/transactions \
-H "Content-Type: application/json" \
-d '{"amount":450,"category":"еда","description":"ланч","date":"2025-09-10"}'
```

## Как запустить тесты

Для запуска тестов и генерации отчета о покрытии используйте скрипт:
```bash
./run_tests.sh
```

*Перед первым запуском сделайте скрипт исполняемым: `chmod +x run_tests.sh`.*
*На Windows используйте терминал Git Bash.*

**Покрытие кода тестами: ~80%**