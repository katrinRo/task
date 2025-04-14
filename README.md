# Todo List API 
Простое приложение для управления списками дел, написанное на Go. Оно позволяет создавать, получать, редактировать и удалять задачи. Проект использует PostgreSQL в качестве базы данных и предоставляет Swagger UI для тестирования API.

Инструкции по установке и запуску.
Клонирование репозитория:

bash
`git clone https://github.com/katrinRo/task.git`
Сборка и запуск контейнеров Docker:

bash
`docker-compose up -d`

Доступ к Swagger UI:
Откройте в браузере http://localhost:8080/swagger/index.html для тестирования API.

Swagger UI не открывается: Убедитесь, что контейнер Swagger UI запущен и доступен по адресу http://localhost:8080.

###Зависимости:
- Go
- PostgreSQL
- Docker 
- Swagger UI
