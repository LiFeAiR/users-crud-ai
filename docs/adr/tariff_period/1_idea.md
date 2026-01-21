## Варианты реализации функции "Тарифы предоставляется пользователю или организации на ограниченный срок от, до"

### Вариант 1: Прямое добавление полей в существующие модели

**Описание**: Добавить поля даты начала и окончания действия тарифа непосредственно в модели User и Organization.

**Преимущества**:
- Простая реализация
- Минимальные изменения в архитектуре
- Прямой доступ к данным без дополнительных запросов

**Недостатки**:
- Нарушает принцип единственной ответственности (в модели будут данные не только о пользователе/организации, но и о времени действия тарифа)
- Не гибкий подход к различным типам ограничений

**Структура изменений**:
```go
// models/user.go
type User struct {
    ID           int           `json:"id"`
    Name         string        `json:"name"`
    Email        string        `json:"email"`
    PasswordHash string        `json:"password_hash"`
    Organization *Organization `json:"organization"`
    TariffID     *int          `json:"tariff_id,omitempty"`
    TariffStart  *time.Time    `json:"tariff_start,omitempty"`  // Новое поле
    TariffEnd    *time.Time    `json:"tariff_end,omitempty"`    // Новое поле
}

// models/organization.go
type Organization struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    TariffID *int   `json:"tariff_id,omitempty"`
    TariffStart *time.Time `json:"tariff_start,omitempty"`  // Новое поле
    TariffEnd   *time.Time `json:"tariff_end,omitempty"`    // Новое поле
}
```

### Вариант 2: Отдельная сущность для хранения временных ограничений

**Описание**: Создать отдельную сущность TariffAssignment с полями дат начала и окончания, которая будет связана с пользователем/организацией и тарифом.

**Преимущества**:
- Четкое разделение ответственности
- Гибкость при работе с различными типами ограничений
- Возможность хранить историю назначений тарифов
- Легко расширяемая архитектура

**Недостатки**:
- Более сложная реализация
- Необходимость дополнительных JOIN-запросов
- Более сложная логика проверки активности тарифа

**Структура изменений**:
```go
// models/tariff_assignment.go
type TariffAssignment struct {
    ID          int        `json:"id"`
    UserID      *int       `json:"user_id,omitempty"`
    OrgID       *int       `json:"org_id,omitempty"`
    TariffID    int        `json:"tariff_id"`
    StartDate   time.Time  `json:"start_date"`
    EndDate     time.Time  `json:"end_date"`
}

// Модифицированные модели
type User struct {
    ID           int           `json:"id"`
    Name         string        `json:"name"`
    Email        string        `json:"email"`
    PasswordHash string        `json:"password_hash"`
    Organization *Organization `json:"organization"`
    TariffID     *int          `json:"tariff_id,omitempty"`
}

type Organization struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    TariffID *int   `json:"tariff_id,omitempty"`
}
```

### Вариант 3: Расширение существующих полей с использованием JSON

**Описание**: Сохранять временные ограничения в JSON-поле внутри существующих моделей.

**Преимущества**:
- Минимальные изменения в структуре БД
- Гибкость в структуре данных
- Простота реализации

**Недостатки**:
- Ограниченная возможность индексации и поиска
- Сложнее валидировать данные
- Потеря типобезопасности

**Структура изменений**:
```go
// models/user.go
type User struct {
    ID           int           `json:"id"`
    Name         string        `json:"name"`
    Email        string        `json:"email"`
    PasswordHash string        `json:"password_hash"`
    Organization *Organization `json:"organization"`
    TariffID     *int          `json:"tariff_id,omitempty"`
    TariffPeriod *TariffPeriod `json:"tariff_period,omitempty"`  // JSON-поле
}

type TariffPeriod struct {
    Start time.Time `json:"start"`
    End   time.Time `json:"end"`
}
```

## Рекомендация

Наиболее подходящим вариантом является **Вариант 1** (прямое добавление полей в существующие модели), поскольку:

1. Он обеспечивает минимальные изменения в архитектуре
2. Простота реализации и поддержки
3. Прямой доступ к данным без дополнительных JOIN-запросов
4. Соответствует текущему стилю проекта
5. Для большинства случаев достаточно простого временного ограничения

Этот подход будет наиболее эффективным для текущей задачи, так как:
- Не требует значительных изменений в существующей логике
- Обеспечивает хорошую производительность
- Прост в тестировании и отладке
- Сохраняет совместимость с текущей архитектурой

Давайте теперь подготовим техническое описание для реализации выбранного варианта:

