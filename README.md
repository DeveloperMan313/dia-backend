# Let There Be Light API

REST API веб-приложения расчета освещения комнат

## Аутентификация

Проект использует **JWT (JSON Web Tokens)** для аутентификации и **ролевую модель** для авторизации.

### Роли пользователей

- **User (0)** - обычные пользователи, могут создавать и управлять своими заявками
- **Moderator (1)** - модераторы, могут управлять лампами и модерировать заявки

## Swagger API Documentation

Интерактивная документация API доступна через Swagger UI:

```
http://localhost:8001/swagger/index.html
```

Генерация новой документации:

```bash
swag init -g cmd/dia-backend/main.go -o docs
```

## Домены API

### Домен лампы (Lamp)

- `GET /api/lamps` - список ламп с фильтрацией по названию
- `GET /api/lamps/:id` - получение лампы по ID
- `POST /api/lamps` - создание новой лампы
- `PUT /api/lamps/:id` - обновление лампы
- `DELETE /api/lamps/:id` - удаление лампы
- `POST /api/lamps/:id/image` - добавление/замена изображения лампы
- `POST /api/lamps/:id/draft` - добавление лампы в черновую заявку

### Домен заявки (LightRequest)

- `GET /api/light-requests/cart` - информация о корзине (ID черновой заявки и количество ламп)
- `GET /api/light-requests` - список заявок с фильтрацией по статусу и дате формирования
- `GET /api/light-requests/:id` - получение заявки с лампами по ID
- `PUT /api/light-requests/:id` - обновление полей заявки
- `PUT /api/light-requests/:id/form` - формирование заявки создателем
- `PUT /api/light-requests/:id/resolve` - разрешение заявки модератором
- `PUT /api/light-requests/:id/reject` - отклонение заявки модератором
- `DELETE /api/light-requests/:id` - удаление заявки

### Домен связь заявка-лампа (LightRequest-Lamp)

- `DELETE /api/light-request-lamps` - удаление лампы из заявки
- `PUT /api/light-request-lamps` - изменение количества/площади лампы в заявке

### Домен пользователь (User)

- `POST /api/users/register` - регистрация пользователя (публичный)
- `POST /api/users/login` - аутентификация (публичный)
- `POST /api/users/logout` - выход из системы (требует JWT)
- `GET /api/users/profile` - получение профиля пользователя (требует JWT)
- `PUT /api/users/profile` - обновление профиля пользователя (требует JWT)

## Модели данных (как они возвращаются API)

### Лампа (Lamp)

```json
{
  "id": 1,
  "title": "Светодиодная лампа",
  "power_w": 10.5,
  "luminous_flux_lm": 800.0,
  "scattering_angle_deg": 120.0,
  "image_url": "http://localhost:9000/lamp-images/lamp_1_1234567890.jpg"
}
```

### Заявка (LightRequest)

```json
{
  "status": 3,
  "max_total_power_w": 1000.0,
  "created_at": "2024-01-15T10:30:00Z",
  "formed_at": "2024-01-16T14:20:00Z",
  "closed_at": "2024-01-17T09:15:00Z",
  "light_request_to_lamp": [
    {
      "request_id": 1,
      "lamp_id": 1,
      "area_m2": 15.0,
      "number": 2,
      "lamp": {
        "id": 1,
        "title": "Светодиодная лампа",
        "power_w": 10.5,
        "luminous_flux_lm": 800.0,
        "scattering_angle_deg": 120.0
      }
    }
  ]
}
```

### Связь заявка-лампа (LightRequestToLamp)

```json
{
  "area_m2": 25.5,
  "number": 3,
  "lamp": {
    "id": 5,
    "title": "Светодиодная лампа Premium",
    "power_w": 15.0,
    "luminous_flux_lm": 1200.0,
    "scattering_angle_deg": 120.0,
    "image_url": "http://localhost:9000/lamp-images/lamp_5.jpg"
  }
}
```

### Статусы заявок

- `1` - Черновик (Draft)
- `2` - Удален (Deleted)
- `3` - Ожидание (Pending)
- `4` - Завершено (Resolved)
- `5` - Отклонено (Rejected)

### Пользователь (User)

```json
{
  "username": "john_doe",
  "role": 0
}
```

**Поле `role`:**

- `0` - User (обычный пользователь)
- `1` - Moderator (модератор)

## Бизнес-логика

### Расчет количества ламп

При завершении заявки рассчитывается количество ламп по формуле:

```
N = (E * S) / Phi
```

Где:

- `N` - количество светильников
- `E` - требуемая освещенность (500 Люкс - стандартное значение для офисных помещений)
- `S` - площадь помещения (м²)
- `Phi` - световой поток одной лампы (Люмен)

### Workflow статусов

- Пользователь: создает черновик → формирует заявку
- Модератор: отклоняет или завершает сформированную заявку
- Системные поля вычисляются автоматически на бэкенде
- Пустая заявка-черновик создается автоматически при добавлении первой лампы

## Конфигурация

### Настройки через переменные окружения:

Расположение env файла:

```
deploy/.env
```

Пример env файла:

```
deploy/.env.example
```

## Запуск

```bash
cd deploy
docker compose up -d
cd ..
go run ./cmd/migrate/
go run ./cmd/dia-backend/
```
