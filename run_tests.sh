#!/bin/bash

# Скрипт для запуска тестов, генерации отчета о покрытии и его открытия.
# Запускать из корневой директории проекта.

# --- Пояснения к командам ---
set -e

# Выводим сообщение о начале процесса для наглядности.
echo "Running tests and generating coverage report..."

# --- Шаг 1: Запуск тестов и создание профиля покрытия ---
go test -v ./ledger/... ./gateway/... -cover -coverprofile=cover.out

echo "-------------------------------------"
echo "Tests passed successfully!"
echo "-------------------------------------"


# --- Шаг 2: Генерация HTML-отчета из профиля ---
go tool cover -html cover.out -o cover.html

echo "Coverage report generated: cover.html"
echo "-------------------------------------"


# --- Шаг 3 (Опционально): Открытие отчета в браузере ---
echo "Attempting to open the report in your default browser..."

# Проверяем, какая ОС используется, и выбираем соответствующую команду.
case "$(uname -s)" in
   Linux*)    xdg-open cover.html ;;
   Darwin*)   open cover.html ;;
   CYGWIN*|MINGW*|MSYS*|windows*) start cover.html ;; # Для Windows (Git Bash, WSL и т.д.)
   *)         echo "Could not detect OS to open browser. Please open cover.html manually." ;;
esac

echo "Script finished."