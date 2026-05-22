# Subscription Service API

REST API сервис для агрегации и управления онлайн-подписками пользователей.  
Тестовое задание для позиции Junior Golang Developer в Effective-mobile


## Subscriptions

| Method | Endpoint | Description |
|---|---|---|
| POST | `/subscriptions/` | Создать подписку |
| GET | `/subscriptions/` | Получить список подписок |
| GET | `/subscriptions/:id` | Получить подписку по ID |
| PUT | `/subscriptions/:id` | Обновить подписку |
| DELETE | `/subscriptions/:id` | Удалить подписку |
| GET | `/subscriptions/total` | Подсчитать общую стоимость |

# Запуск и т.п

docker compose up --build

# Swagger

Swagger документация будет доступна по адресу:
http://localhost:8080/swagger/index.html
