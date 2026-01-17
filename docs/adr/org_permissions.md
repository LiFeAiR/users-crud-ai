## Что было реализовано:

Я успешно реализовал методы AddOrganizationPermissions и DeleteOrganizationPermissions, описанные в файле api/grpc/api.proto.

Были выполнены следующие шаги:

1. Добавлены методы работы с правами в интерфейс OrganizationRepository:
    - GetOrganizationPermissions - получение прав организации
    - AddOrganizationPermissions - добавление прав к организации
    - DeleteOrganizationPermissions - удаление прав из организации

2. Реализованы методы в organization_repository.go:
    - Добавлена таблица связи organization_permissions в БД
    - Реализована логика работы с правами организаций
    - Добавлена проверка существования организации и прав перед операциями

3. Созданы обработчики для gRPC методов:
    - internal/handlers/organization_add_permissions.go - обработчик AddOrganizationPermissions
    - internal/handlers/organization_delete_permissions.go - обработчик DeleteOrganizationPermissions

Все методы следуют единому шаблону и интегрированы в существующую архитектуру проекта. Методы полностью соответствуют описанию в api.proto и готовы к использованию.

Я успешно реализовал методы AddOrganizationRoles и DeleteOrganizationRoles, описанные в файле api/grpc/api.proto.

Выполненные шаги:

1. Добавлены методы работы с ролями в интерфейс OrganizationRepository:
   - GetOrganizationRoles - получение ролей организации
   - AddOrganizationRoles - добавление ролей к организации
   - DeleteOrganizationRoles - удаление ролей из организации

2. Реализованы методы в organization_repository.go:
   - Добавлена таблица связи organization_roles в БД
   - Реализована логика работы с ролями организаций
   - Добавлена проверка существования организации и ролей перед операциями

3. Созданы обработчики для gRPC методов:
   - internal/handlers/organization_add_roles.go - обработчик AddOrganizationRoles
   - internal/handlers/organization_delete_roles.go - обработчик DeleteOrganizationRoles
4. Обновлен метод GetOrganization (стоит подумать над упрощением запроса)
```sql
SELECT DISTINCT p.id, p.name, p.code, p.description
FROM permissions p
        LEFT JOIN organization_permissions op ON p.id = op.permission_id and op.organization_id = $1
        LEFT JOIN organization_roles ro ON ro.organization_id = $1
        LEFT JOIN role_permissions rp ON rp.role_id = ro.role_id and p.id = op.permission_id
WHERE ro.organization_id = $1
   or op.organization_id = $1
ORDER BY p.id
```

Все методы следуют единому шаблону и интегрированы в существующую архитектуру проекта. Методы полностью соответствуют описанию в api.proto и готовы к использованию.