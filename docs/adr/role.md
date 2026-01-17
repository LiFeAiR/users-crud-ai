1. Создана модель Role в `internal/models/role.go` с полем для связи с Permission
2. Добавлен интерфейс и реализация репозитория для Role в `internal/repository/role_repository.go`
3. Созданы хендлеры для всех CRUD операций с Role:
    - `internal/handlers/role_create.go`
    - `internal/handlers/role_get.go`
    - `internal/handlers/role_update.go`
    - `internal/handlers/role_delete.go`
    - `internal/handlers/roles_get.go`
4. Добавлены методы управления правами ролей:
    - `internal/handlers/role_add_permissions.go`
    - `internal/handlers/role_delete_permissions.go`
5. Обновлен протокол API в `api/grpc/api.proto` для поддержки операций с Role
6. Обновлен базовый обработчик в `internal/handlers/base_handler.go` для поддержки Role
7. Обновлена инициализация сервера в `internal/server/server.go` для поддержки Role
8. Созданы тесты для новых компонентов:
    - `internal/repository/role_repository_test.go`
    - `internal/handlers/role_create_test.go`
    - `internal/handlers/mock_role_repository.go`

Все компоненты соответствуют существующей архитектуре проекта и следуют тем же паттернам, что и другие сущности (User, Organization, Permission). Связь один ко многим между Role и Permission реализована через таблицу связи role_permissions, что позволяет одной роли иметь несколько прав и одному праву принадлежать к нескольким ролям.

API теперь поддерживает следующие операции с ролями:
- Создание роли (POST /api/roles)
- Получение роли по ID (GET /api/role/{id})
- Обновление роли (PUT /api/role)
- Удаление роли (DELETE /api/role/{id})
- Получение списка ролей (GET /api/roles)
- Добавление прав к роли (POST /api/role/{role_id}/permissions)
- Удаление прав из роли (DELETE /api/role/{role_id}/permissions)