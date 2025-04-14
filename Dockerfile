# Используем базовый образ Go
FROM golang:1.24.0-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Обновляем зависимости
RUN go mod tidy

# Копируем весь исходный код
COPY . .

# Скомпилируйте приложение
RUN go build -o task .

# Используем легковесный образ для запуска приложения
FROM alpine

# Копируем скомпилированное приложение
COPY --from=builder /app/task .

COPY ./sql sql

COPY ./docs docs


# Указываем команду для запуска приложения
CMD ["./task"]
