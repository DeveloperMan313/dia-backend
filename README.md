# Let There Be Light API

REST API веб-приложения расчета освещения комнат

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
- `PUT /api/light-requests/:id/resolve` - завершение заявки модератором
- `PUT /api/light-requests/:id/reject` - отклонение заявки модератором
- `DELETE /api/light-requests/:id` - удаление заявки

### Домен связь заявка-лампа (LightRequest-Lamp)

- `DELETE /api/light-request-lamps` - удаление лампы из заявки
- `PUT /api/light-request-lamps` - изменение количества/площади лампы в заявке

### Домен пользователь (User)

- `POST /api/users/register` - регистрация пользователя
- `GET /api/users/profile` - получение профиля пользователя
- `PUT /api/users/profile` - обновление профиля пользователя
- `POST /api/users/login` - аутентификация
- `POST /api/users/logout` - деавторизация

## Модели данных

### Лампа (Lamp)

```json
{
  "id": 1,
  "title": "Светодиодная лампа",
  "power_w": 10.5,
  "luminous_flux_lm": 800.0,
  "scattering_angle_deg": 120.0,
  "is_deleted": false,
  "image_url": "http://localhost:9000/lamp-images/lamp_1_1234567890.jpg"
}
```

### Заявка (LightRequest)

```json
{
  "id": 1,
  "status": 3,
  "user_id": 1,
  "moderator_id": 2,
  "max_total_power_w": 1000.0,
  "created_at": "2024-01-15T10:30:00Z",
  "formed_at": "2024-01-16T14:20:00Z",
  "closed_at": "2024-01-17T09:15:00Z",
  "user": {
    "id": 1,
    "username": "creator"
  },
  "moderator": {
    "id": 2,
    "username": "moderator"
  },
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

### Статусы заявок

- `1` - Черновик (Draft)
- `2` - Удален (Deleted)
- `3` - Ожидание (Pending)
- `4` - Завершено (Resolved)
- `5` - Отклонено (Rejected)

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

**Примечание**: Расчетные данные (общая мощность, дата доставки, количество ламп) возвращаются в ответе API, но не сохраняются в базе данных для соблюдения нормализации.

**Создатель**: ID пользователя зафиксирован как 1.

### Workflow статусов

- Пользователь: создает черновик → формирует заявку
- Модератор: отклоняет или завершает сформированную заявку
- Системные поля вычисляются автоматически на бэкенде
- Пустая заявка-черновик создается автоматически при добавлении первой лампы

## Запуск

```bash
go run cmd/dia-backend/main.go
```

## Конфигурация

Настройки через переменные окружения:

- `LISTEN_PORT` - порт для запуска сервера
- `DATABASE_URL` - строка подключения к PostgreSQL
- `MINIO_ENDPOINT` - endpoint MinIO для хранения изображений
