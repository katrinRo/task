-- Создание таблицы
CREATE TABLE IF NOT EXISTS task (
                          id BIGSERIAL PRIMARY KEY,
                          title VARCHAR(100) NOT NULL,
                          description VARCHAR(400),
                          date DATE NOT NULL,
                          status BOOLEAN NOT NULL
);