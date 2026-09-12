
# Genealogy Tree

Backend MVP генеалогического дерева на Go.

## Запуск

```sh
cp .env.example .env
make env-up
make migrate-up
make run
```

Для полного удаления контейнеров и volumes окружения используй `make env-cleanup` и подтверди операцию символом `y`.

## MinIO и фотографии

Локальный MinIO запускается вместе с PostgreSQL:

```sh
make env-up
```

Адреса локального MinIO:

- S3 API: `http://localhost:9000`
- Web Console: `http://localhost:9001`
- Логин: значение `MINIO_ACCESS_KEY` из `.env`
- Пароль: значение `MINIO_SECRET_KEY` из `.env`
- Bucket: значение `MINIO_BUCKET` из `.env`

При запуске приложения bucket создаётся автоматически, если его ещё нет. Файлы хранятся локально в `out/minio`, поэтому удаление контейнера не удаляет фотографии. Полная очистка через `make env-cleanup` удаляет и этот каталог.

### Загрузка фотографии

Сначала создай person через `POST /api/v1/persons` и возьми его ULID из ответа. Затем отправь запрос:

```http
POST /api/v1/persons/{id}/photo
Content-Type: multipart/form-data
```

В Postman выбери `Body -> form-data`, создай поле `file`, установи для него тип `File` и выбери фотографию. Максимальный размер файла: 10 MB.

Пример через curl:

```sh
curl -X POST http://localhost:8080/api/v1/persons/{id}/photo \
  -F "file=@/path/to/photo.jpg"
```

После успешной загрузки в поле `photo_url` person сохраняется ключ объекта MinIO, например `persons/{id}/photo`.

### Получение фотографии

```http
GET /api/v1/persons/{id}/photo
```

Тело запроса не требуется. API возвращает бинарное содержимое фотографии.

### Удаление фотографии

```http
DELETE /api/v1/persons/{id}/photo
```

Тело запроса не требуется. Объект удаляется из MinIO, а `photo_url` person очищается.

Повторная загрузка фотографии заменяет старый файл и удаляет прежний объект из MinIO.

Все идентификаторы в API имеют формат ULID, например `01JQ2Q4K7Y8F6M2Z3N4P5R6S7T`.

Проверка доступности: `GET http://localhost:8080/health`.

API использует префикс `/api/v1`:

### Persons

`POST /api/v1/persons` - создать человека. Тело обязательно:

```json
{
  "first_name": "Ivan",
  "last_name": "Petrov",
  "birth_date": "1980-01-02",
  "death_date": null,
  "gender": "male",
  "photo_url": null,
  "metadata": {"city": "Moscow"}
}
```

Обязательные поля: `first_name`, `last_name`. Допустимые значения `gender`: `male`, `female`, `other`.

`GET /api/v1/persons` - получить список людей. Тело не требуется.

`GET /api/v1/persons/{id}` - получить человека. Тело не требуется.

`PATCH /api/v1/persons/{id}` - изменить человека. Тело обязательно, все поля необязательны:

```json
{
  "last_name": "Sidorov",
  "death_date": "2024-05-10",
  "metadata": {"city": "Berlin"}
}
```

Для очистки nullable-поля передайте `null`, например `{ "photo_url": null }`.

`DELETE /api/v1/persons/{id}` - удалить человека и связанные relationships. Тело не требуется.

### Relationships

`POST /api/v1/relationships` - создать связь. Тело обязательно:

Для связи родитель-ребёнок:

```json
{
  "person1_id": "01JQ2Q4K7Y8F6M2Z3N4P5R6S7T",
  "person2_id": "01JQ2Q4K8Z8F6M2Z3N4P5R6S7U",
  "type": "parent_child",
  "direction": "parent"
}
```

Для связи супругов поле `direction` не передаётся:

```json
{
  "person1_id": "01JQ2Q4K7Y8F6M2Z3N4P5R6S7T",
  "person2_id": "01JQ2Q4K8Z8F6M2Z3N4P5R6S7U",
  "type": "spouse"
}
```

`GET /api/v1/relationships/{id}` - получить связь. Тело не требуется.

`DELETE /api/v1/relationships/{id}` - удалить связь. Тело не требуется.
