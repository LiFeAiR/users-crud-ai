Я успешно реализовал CRUD-операции для сущности "тариф" с возможностью связи один ко многим с сущностью "роль". Вот что было сделано:

## 1. Создана модель данных для тарифа
- Создан файл `internal/models/tariff.go` с моделью Tariff, содержащей поля: id, name, description, price
- Добавлена связь один ко многим с Role через поле Roles

## 2. Реализован интерфейс репозитория
- Добавлен интерфейс `TariffRepository` в `internal/repository/interface.go` с методами:
   - CreateTariff
   - GetTariffByID
   - UpdateTariff
   - DeleteTariff
   - GetTariffs
   - GetTariffWithRoles
   - AddTariffRoles
   - DeleteTariffRoles
   - InitDB

## 3. Реализован репозиторий
- Создан файл `internal/repository/tariff_repository.go` с реализацией всех методов интерфейса
- Добавлена таблица `tariffs` в БД с полями id, name, description, price
- Добавлена таблица связи `tariff_roles` для реализации связи один ко многим

## 4. Добавлены обработчики
- Созданы файлы обработчиков в `internal/handlers/`:
   - tariff_create.go
   - tariff_get.go
   - tariff_update.go
   - tariff_delete.go
   - tariffs_get.go
- Обновлен `internal/handlers/base_handler.go` для поддержки нового репозитория

## 5. Обновлен protobuf-файл
- Добавлены сообщения для тарифа в `api/grpc/api.proto`:
   - Tariff
   - TariffCreateRequest
   - TariffUpdateRequest
   - TariffsResponse
- Добавлены RPC-методы в сервис CrudService:
   - GetTariffs
   - CreateTariff
   - GetTariff
   - UpdateTariff
   - DeleteTariff

## Особенности реализации:
- Связь один ко многим реализована через таблицу `tariff_roles` (many-to-many)
- Все CRUD-операции поддерживают работу с ролью через связи
- Соблюдены существующие паттерны кодирования проекта
- Все методы имеют соответствующую валидацию и обработку ошибок
- Поддерживается работа с JSON и gRPC интерфейсами

Реализация полностью соответствует архитектуре проекта и позволяет выполнять все необходимые операции с сущностью "тариф" и его связями с "ролями".