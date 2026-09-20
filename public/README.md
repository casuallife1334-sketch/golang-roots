# Roots frontend

React + TypeScript + Vite приложение в `/public`. Оно использует только реальный Go API по относительному адресу `/api/v1`; Vite в development проксирует `/api` на `http://localhost:8080`.

## Запуск

Требуется Node.js 20+ и npm.

```sh
npm install
npm run dev
```

Frontend: `http://localhost:5173`. Backend запускается существующими командами из корня проекта: `make docker-deploy` или `make env-up && make migrate-up && make run`.

## Запуск всей связки в Docker

Frontend собирается multi-stage Dockerfile и запускается через Nginx. Nginx раздаёт SPA и проксирует `/api` к compose-сервису `genealogy`, поэтому frontend-контейнеру не нужен доступ к `localhost:8080`.

Из корня проекта:

```sh
make docker-deploy
```

Либо напрямую:

```sh
docker compose up --build
```

После запуска приложение доступно по адресу `http://localhost:5173`. Порт можно изменить через `FRONTEND_PORT` в `.env`:

```env
FRONTEND_PORT=5173
```

SPA-маршруты (`/login`, `/register`, `/trees/...`, `/settings`) обслуживаются через fallback на `index.html`. Запросы к `/api/v1` проходят внутри compose-сети по адресу `http://genealogy:8080`.

Проверки:

```sh
npm run typecheck
npm run build
npm run lint
npm test
```

Для production задайте публичный base URL через reverse proxy: `/api` должен проксироваться на Go API, а неизвестные SPA-маршруты возвращать `public/index.html`. JWT хранится только в `sessionStorage` текущей вкладки.

## Реализовано

- `/login`, `/register`, проверка `/users/me`, logout и обработка 401/403;
- выбор, создание, переименование и удаление деревьев;
- реальные persons, relationships и фотографии через scoped API;
- вычисляемая страница «Семьи» без отдельного backend endpoint;
- локальные настройки отображения под ключом user ID;
- React Flow-граф с детерминированной раскладкой parent-child/spouse связей;
- пустые, загрузочные и ошибочные состояния.

Референсы из `../docs/design-references` использованы как визуальное направление. Не реализованы функции, которых нет в API: приглашения, editor/viewer управление, timeline, GEDCOM, premium, 2FA, refresh token и смена профиля.

Подтверждённое ограничение backend: nullable-поля в PATCH человека могут не различать явный `null` и отсутствие поля, поэтому интерфейс не обещает очистку уже заполненных дат/пола и не удаляет запись для обхода этого ограничения.
