# MultiSub — Подробный план разработки модулей

**Версия документа:** 1.0  
**Дата создания:** 27 февраля 2026  
**Связанный документ:** [TECHNICAL_SPECIFICATION.md](./TECHNICAL_SPECIFICATION.md)

---

## 📋 Оглавление

1. [Общие принципы разработки](#0-общие-принципы-разработки)
2. [Модуль 1: Auth](#модуль-1-аутентификация-и-авторизация)
3. [Модуль 2: User Profile](#модуль-2-профиль-пользователя)
4. [Модуль 3: Subscriptions](#модуль-3-управление-подписками)
5. [Модуль 4: Bank Integration](#модуль-4-интеграция-с-банками)
6. [Модуль 5: Email Parsing](#модуль-5-парсинг-почты)
7. [Модуль 6: Analytics](#модуль-6-аналитика)
8. [Модуль 7: Usage Tracking](#модуль-7-отслеживание-использования)
9. [Модуль 8: Prediction](#модуль-8-прогнозирование)
10. [Модуль 9: Notifications](#модуль-9-уведомления)
11. [Модуль 10: Catalog & Alternatives](#модуль-10-каталог-сервисов-и-альтернативы)
12. [Модуль 11: Cancellation](#модуль-11-отмена-подписок)
13. [Модуль 12: Admin](#модуль-12-административная-панель)
14. [Модуль 13: Mobile App](#модуль-13-мобильное-приложение)
15. [Модуль 14: Web App](#модуль-14-веб-приложение)

---

## 0. Общие принципы разработки

### 0.1 Workflow разработки

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Планирование│───▶│  Дизайн API │───▶│ Разработка  │───▶│  Code Review│───▶│Тестирование │
│  задач       │    │  и схемы БД │    │             │    │             │    │             │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘    └──────┬──────┘
                                                                                    │
                                                         ┌──────────────────────────▼──────┐
                                                         │         Deploy to Staging        │
                                                         └──────────────────────────┬──────┘
                                                                                    │
                                                         ┌──────────────────────────▼──────┐
                                                         │      QA + Acceptance Tests       │
                                                         └──────────────────────────┬──────┘
                                                                                    │
                                                         ┌──────────────────────────▼──────┐
                                                         │         Deploy to Prod           │
                                                         └─────────────────────────────────┘
```

### 0.2 Definition of Done для каждой задачи

- [ ] Код написан и соответствует code style
- [ ] Unit-тесты покрывают ≥80% строк
- [ ] Integration-тесты написаны для всех API endpoints
- [ ] Код прошёл code review (минимум 1 апрувер)
- [ ] Документация обновлена (Swagger/OpenAPI)
- [ ] Нет критических и высоких уязвимостей (SAST)
- [ ] Миграции БД написаны и обратно совместимы
- [ ] Задача задеплоена на staging и проверена
- [ ] Acceptance criteria выполнены и подтверждены PM

### 0.3 Соглашения по тестированию

| Тип теста | Инструмент (Backend) | Инструмент (Web) | Инструмент (Mobile) | Покрытие |
|-----------|----------------------|-----------------|---------------------|----------|
| Unit | Go testing + testify | Jest + RTL | Jest + RN Testing Library | ≥80% |
| Integration | httptest + testcontainers | MSW + Jest | Detox (E2E) | Ключевые flows |
| E2E | Playwright (API) | Playwright | Detox | Critical paths |
| Load | k6 | - | - | SLA метрики |
| Security | gosec + OWASP ZAP | - | - | Перед релизом |

### 0.4 Структура папок Backend (Go)

```
internal/
  auth/
    handler.go          # HTTP обработчики
    service.go          # Бизнес-логика
    repository.go       # Работа с БД
    dto.go              # Data Transfer Objects
    handler_test.go
    service_test.go
    repository_test.go
  subscriptions/
    ...
pkg/
  middleware/
  errors/
  validator/
  logger/
migrations/
  001_create_users.sql
  002_create_subscriptions.sql
docs/
  swagger.yaml
```

### 0.5 Среды окружения

| Среда | Назначение | URL |
|-------|-----------|-----|
| `local` | Локальная разработка | localhost |
| `dev` | Ветка develop, автодеплой | dev.multisub.app |
| `staging` | Pre-prod, QA тестирование | staging.multisub.app |
| `production` | Боевой сервер | multisub.app |

---

## Принципы проектирования фронтенд-части

### FE-1. Ключевые правила — обязательны к соблюдению

| Правило | Описание |
|---------|----------|
| **Запрет эмодзи** | Нигде в UI не используются символы эмодзи (ни в кнопках, ни в заголовках, ни в тостах, ни в MD-файлах, отображаемых пользователю) |
| **Только react-icons** | Все иконки берутся из библиотеки `react-icons`. Никакого инлайнового SVG, никаких других icon-библиотек (Heroicons, Lucide, Material Icons) — только `react-icons/fi`, `react-icons/ri`, `react-icons/md`, `react-icons/tb` |
| **CSS-переменные** | Все цвета, отступы, радиусы — через Tailwind-токены или CSS custom properties. Никаких хардкодных hex-значений в JSX |
| **Design tokens first** | При добавлении нового цвета/шрифта/радиуса сначала добавляем токен в `tailwind.config.ts`, потом используем |

### FE-2. Структура токенов дизайн-системы

```ts
// tailwind.config.ts — единый источник правды
export default {
  theme: {
    extend: {
      colors: {
        brand: {
          50:  'var(--brand-50)',
          500: 'var(--brand-500)',  // основной акцент — нейтральный тёмный/синий
          900: 'var(--brand-900)',
        },
        surface: {
          DEFAULT: 'var(--surface)',      // фон карточек
          raised:  'var(--surface-raised)',
          overlay: 'var(--surface-overlay)',
        },
        text: {
          primary:   'var(--text-primary)',
          secondary: 'var(--text-secondary)',
          muted:     'var(--text-muted)',
        },
        status: {
          success: 'var(--status-success)',
          warning: 'var(--status-warning)',
          danger:  'var(--status-danger)',
          info:    'var(--status-info)',
        },
        border: 'var(--border)',
      },
      fontFamily: {
        sans: ['var(--font-sans)', 'system-ui', 'sans-serif'],
        mono: ['var(--font-mono)', 'monospace'],
      },
      borderRadius: {
        card: '12px',
        pill: '9999px',
      },
    },
  },
}
```

```css
/* globals.css — svetlая тема (дефолт) */
:root {
  --brand-50:  #f0f4ff;
  --brand-500: #3b5bdb;
  --brand-900: #1a2a6c;

  --surface:         #ffffff;
  --surface-raised:  #f8f9fa;
  --surface-overlay: #e9ecef;

  --text-primary:   #1a1a2e;
  --text-secondary: #4a4a6a;
  --text-muted:     #9a9ab3;

  --status-success: #2f9e44;
  --status-warning: #f59f00;
  --status-danger:  #e03131;
  --status-info:    #1971c2;

  --border: #dee2e6;
}

/* Тёмная тема */
[data-theme="dark"] {
  --surface:         #0d1117;
  --surface-raised:  #161b22;
  --surface-overlay: #21262d;
  --text-primary:   #e6edf3;
  --text-secondary: #8b949e;
  --text-muted:     #484f58;
  --border:         #30363d;
}
```

### FE-3. Использование иконок (`react-icons`)

```tsx
// ПРАВИЛЬНО — единый стиль через Feather Icons (react-icons/fi)
import { FiHome, FiCreditCard, FiBell, FiSettings, FiLogOut, FiPlus } from 'react-icons/fi'
import { FiTrendingUp, FiCalendar, FiAlertCircle } from 'react-icons/fi'

// Для специфичных иконок — Remix Icons (react-icons/ri)
import { RiSubscriptLine, RiBankLine } from 'react-icons/ri'

// Обёртка для унифицированного размера
type IconProps = { size?: number; className?: string }

export function NavIcon({ icon: Icon, size = 20, className }: IconProps & { icon: IconType }) {
  return <Icon size={size} className={className} aria-hidden="true" />
}

// Использование
<NavIcon icon={FiHome} size={20} className="text-text-secondary" />

// ЗАПРЕЩЕНО — никаких эмодзи
// return <span>📊 Аналитика</span>       ← НЕЛЬЗЯ
// return <svg>...</svg>                  ← нельзя инлайн без крайней нужды
// import { HomeIcon } from '@heroicons/react/24/outline'  ← нельзя, только react-icons
```

**Используемые наборы из `react-icons`:**

| Набор | Префикс импорта | Применение |
|-------|----------------|------------|
| Feather Icons | `react-icons/fi` | Основная навигация, действия, UI |
| Remix Icons | `react-icons/ri` | Финансы, подписки, специфика |
| Material Design | `react-icons/md` | Системные состояния, статусы |
| Tabler Icons | `react-icons/tb` | Дополнительные, если нет в fi/ri |

### FE-4. Архитектура компонентов (Next.js, App Router)

```
apps/web/
  app/
    (auth)/          layout.tsx — AuthLayout (без sidebar)
      login/
      register/
    (dashboard)/     layout.tsx — DashboardLayout с sidebar
      page.tsx       — /  =>  Dashboard
      subscriptions/ — /subscriptions
      analytics/     — /analytics
      settings/      — /settings
  components/
    ui/              — атомарные компоненты (Button, Badge, Card, Input, Modal)
    layout/          — Sidebar, Header, PageHeader, Breadcrumbs
    features/        — feature-specific (SubscriptionCard, CostChart, NotificationList)
    icons/           — обёртки над react-icons (если нужны кастомные размеры)
  lib/
    api.ts           — TanStack Query ключи + fetchers
    store.ts         — Zustand глобальный стор
    utils.ts         — cn(), formatCurrency(), formatDate()
  styles/
    globals.css      — CSS variables + base styles
```

**Структура компонента (образец):**

```tsx
// components/ui/Button.tsx
import { type ButtonHTMLAttributes, forwardRef } from 'react'
import { FiLoader } from 'react-icons/fi'
import { cn } from '@/lib/utils'

type Variant = 'primary' | 'secondary' | 'ghost' | 'danger'
type Size    = 'sm' | 'md' | 'lg'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant
  size?: Size
  loading?: boolean
  leftIcon?: React.ReactNode
}

const variantClasses: Record<Variant, string> = {
  primary:   'bg-brand-500 text-white hover:bg-brand-600 active:bg-brand-700',
  secondary: 'bg-surface-raised text-text-primary border border-border hover:bg-surface-overlay',
  ghost:     'text-text-secondary hover:bg-surface-raised hover:text-text-primary',
  danger:    'bg-status-danger text-white hover:opacity-90',
}

const sizeClasses: Record<Size, string> = {
  sm: 'h-8  px-3 text-sm  gap-1.5',
  md: 'h-10 px-4 text-sm  gap-2',
  lg: 'h-12 px-6 text-base gap-2',
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ variant = 'primary', size = 'md', loading, leftIcon, className, children, disabled, ...props }, ref) => (
    <button
      ref={ref}
      disabled={disabled || loading}
      className={cn(
        'inline-flex items-center justify-center rounded-card font-medium',
        'transition-all duration-150 disabled:pointer-events-none disabled:opacity-50',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500 focus-visible:ring-offset-2',
        variantClasses[variant],
        sizeClasses[size],
        className,
      )}
      {...props}
    >
      {loading ? <FiLoader size={16} className="animate-spin" aria-hidden="true" /> : leftIcon}
      {children}
    </button>
  ),
)
Button.displayName = 'Button'
```

### FE-5. Архитектура компонентов (React Native, Expo)

```
apps/mobile/
  app/
    (auth)/          — AuthStack (без tab bar)
      login.tsx
      register.tsx
    (tabs)/          — TabNavigator
      index.tsx      — Dashboard
      subscriptions.tsx
      analytics.tsx
      settings.tsx
  components/
    ui/              — атомарные (Button, Card, Badge, Input, BottomSheet)
    layout/          — Screen, SafeAreaWrapper, TabBar
    features/        — SubscriptionCard, CostChart, NotificationItem
  lib/
    api.ts
    store.ts
    utils.ts
  hooks/             — useSubscriptions, useAnalytics, useNotifications
```

**Правила для React Native:**

```tsx
// ПРАВИЛЬНО: иконки через react-icons/fi (expo с vector icons)
// В React Native используем @expo/vector-icons как обёртку,
// но именуем константы в стиле react-icons для единообразия
import { Feather } from '@expo/vector-icons'

// Унифицированная обёртка — имитирует react-icons API
export function Icon({ name, size = 20, color }: { name: string; size?: number; color?: string }) {
  return <Feather name={name as any} size={size} color={color ?? colors.text.secondary} />
}

// Использование — аналогично веб
<Icon name="home" size={22} />
<Icon name="credit-card" size={22} />
<Icon name="bell" size={22} />

// НИКАКИХ эмодзи в Text-компонентах
// <Text>📊 Расходы</Text>  ← ЗАПРЕЩЕНО
// <Text>Расходы</Text>     ← правильно, с Icon рядом
```

### FE-6. Соглашения по именованию

| Сущность | Формат | Пример |
|----------|--------|--------|
| Компоненты | PascalCase | `SubscriptionCard`, `CostChart` |
| Хуки | camelCase с `use` | `useSubscriptions`, `useCurrency` |
| Утилиты | camelCase | `formatCurrency`, `truncateText` |
| Константы | UPPER_SNAKE | `MAX_RETRY_COUNT`, `DEFAULT_CURRENCY` |
| CSS-классы | kebab-case (если кастомные) | `subscription-card--active` |
| Файлы компонентов | PascalCase | `Button.tsx`, `SubscriptionCard.tsx` |
| Файлы утилит/хуков | camelCase | `formatCurrency.ts`, `useSubscriptions.ts` |

### FE-7. Типизация и состояние

```ts
// Все API-ответы имеют общий wrapper
interface ApiResponse<T> {
  data: T
  error?: string
  meta?: { total: number; page: number; pageSize: number }
}

// Zustand store — одна ответственность на стор
interface SubscriptionStore {
  filter: SubscriptionFilter
  setFilter: (f: Partial<SubscriptionFilter>) => void
  // НЕТ данных с сервера — только UI-состояние
  // данные — через TanStack Query
}

// TanStack Query ключи — фабрика
export const subscriptionKeys = {
  all:    () => ['subscriptions'] as const,
  list:   (f: SubscriptionFilter) => [...subscriptionKeys.all(), 'list', f] as const,
  detail: (id: string) => [...subscriptionKeys.all(), 'detail', id] as const,
}
```

### FE-8. Запрещённые паттерны

| Паттерн | Причина | Альтернатива |
|---------|---------|--------------|
| Эмодзи в UI-тексте | Непоследовательность, плохая доступность | react-icons |
| Сторонние icon-библиотеки | Дублирование зависимостей | react-icons |
| Инлайн-стили `style={{}}` | Нет design tokens | Tailwind классы |
| Хардкодные цвета `#3b5bdb` | Нет тёмной темы | `text-brand-500` |
| `any` в TypeScript | Теряется безопасность типов | Явные типы / `unknown` |
| Мутация состояния напрямую | Непредсказуемое поведение | Zustand `set()` |
| `useEffect` для выборки данных | Водопад запросов, нет кеша | TanStack Query |
| God-компоненты (>200 строк) | Сложность поддержки | Декомпозиция |

---

## Модуль 1: Аутентификация и авторизация

### Фаза 1: Планирование задач

**Длительность:** 2 дня  
**Участники:** Tech Lead, Backend Developer, UI/UX Designer

#### 1.1 Задачи планирования

- [ ] Проанализировать требования к безопасности (OWASP Top 10)
- [ ] Выбрать стратегию токенов: `access (15 мин) + refresh (7 дней)`
- [ ] Определить OAuth-провайдеры: Google, Яндекс, VK ID
- [ ] Составить схему хранения пароля: `bcrypt cost=12`
- [ ] Согласовать формат JWT payload:
  ```json
  {
    "sub": "user_uuid",
    "email": "user@example.com",
    "role": "user",
    "iat": 1700000000,
    "exp": 1700000900
  }
  ```
- [ ] Утвердить API контракты (OpenAPI spec)
- [ ] Создать тикеты в Jira/Linear для всех задач AUTH-001..012
- [ ] Оценить риски: MITM, brute-force, token leakage

#### 1.2 Схема базы данных

```sql
-- Таблицы: users, refresh_tokens, email_verifications, oauth_accounts, sessions

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,  -- bcrypt hash токена
    device_info JSONB,                 -- {device, os, browser, ip}
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE email_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE oauth_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,      -- google, yandex, vk
    provider_user_id VARCHAR(255) NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, provider_user_id)
);

CREATE TABLE login_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255),
    ip_address INET,
    success BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_login_attempts_email_created ON login_attempts(email, created_at);
CREATE INDEX idx_login_attempts_ip_created ON login_attempts(ip_address, created_at);
```

---

### Фаза 2: Разработка

**Длительность:** ~4 недели

#### 2.1 AUTH-001: Регистрация (email/пароль)

**Backend (Go):**

```go
// internal/auth/dto.go
type RegisterRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8,max=72"`
    FirstName string `json:"first_name" validate:"required,min=1,max=100"`
    LastName  string `json:"last_name" validate:"required,min=1,max=100"`
}

type RegisterResponse struct {
    UserID  string `json:"user_id"`
    Message string `json:"message"` // "Verification email sent"
}

// internal/auth/service.go
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
    // 1. Проверить, что email не занят
    // 2. Захешировать пароль bcrypt(cost=12)
    // 3. Создать пользователя в транзакции
    // 4. Создать код верификации (6 цифр, TTL=24h)
    // 5. Отправить email через очередь (async)
    // 6. Вернуть user_id
}
```

**Шаги разработки:**
1. Написать DTO `RegisterRequest` с валидаторами
2. Написать `UserRepository.Create()`
3. Написать `AuthService.Register()` с транзакцией
4. Написать HTTP handler `POST /api/v1/auth/register`
5. Настроить rate limiting: 5 req/min per IP
6. Добавить в Swagger

#### 2.2 AUTH-003: JWT авторизация

**Алгоритм логина:**
```
1. Проверить email существует → 401 если нет
2. Проверить bcrypt(password, hash) → 401 если не совпадает
3. Проверить email верифицирован → 403 если нет
4. Проверить brute-force: < 5 попыток за последние 15 мин → 429 если превышено
5. Создать access JWT (RS256, 15 мин)
6. Создать refresh token (random 32 bytes, сохранить bcrypt hash в БД)
7. Записать успешную попытку входа
8. Вернуть {access_token, refresh_token, expires_in}
```

#### 2.3 AUTH-009: Двухфакторная аутентификация (TOTP)

**Библиотека:** `github.com/pquerna/otp`

```go
// Генерация TOTP секрета
key, _ := totp.Generate(totp.GenerateOpts{
    Issuer:      "MultiSub",
    AccountName: user.Email,
})

// Верификация кода
valid := totp.Validate(code, key.Secret())
```

**Шаги включения 2FA:**
1. `POST /auth/2fa/enable` → генерировать секрет, вернуть QR-код
2. Пользователь сканирует в Google Authenticator
3. `POST /auth/2fa/verify` с первым кодом → сохранить секрет в БД
4. Хранить `backup_codes` (10 одноразовых кодов) зашифрованными

---

### Фаза 3: Написание тестов

**Длительность:** 1 неделя

#### 3.1 Unit-тесты (Go)

```go
// internal/auth/service_test.go

func TestRegister_Success(t *testing.T) {
    mockRepo := mocks.NewUserRepository(t)
    mockEmailSvc := mocks.NewEmailService(t)
    svc := NewAuthService(mockRepo, mockEmailSvc)

    mockRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, ErrNotFound)
    mockRepo.On("Create", ctx, mock.AnythingOfType("*User")).Return(testUser, nil)
    mockEmailSvc.On("SendVerification", ctx, mock.Anything).Return(nil)

    resp, err := svc.Register(ctx, RegisterRequest{
        Email:     "test@example.com",
        Password:  "SecurePass123!",
        FirstName: "Иван",
        LastName:  "Иванов",
    })

    require.NoError(t, err)
    assert.NotEmpty(t, resp.UserID)
    mockRepo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) { ... }
func TestRegister_WeakPassword(t *testing.T) { ... }
func TestLogin_Success(t *testing.T) { ... }
func TestLogin_WrongPassword(t *testing.T) { ... }
func TestLogin_BruteForceBlocked(t *testing.T) { ... }
func TestRefreshToken_Success(t *testing.T) { ... }
func TestRefreshToken_Expired(t *testing.T) { ... }
func TestRefreshToken_Revoked(t *testing.T) { ... }
```

#### 3.2 Integration-тесты (httptest)

```go
// internal/auth/handler_test.go

func TestRegisterEndpoint(t *testing.T) {
    app := setupTestApp(t)  // реальная БД через testcontainers

    body := `{"email":"test@e.com","password":"Pass123!","first_name":"A","last_name":"B"}`
    req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, _ := app.Test(req)
    assert.Equal(t, 201, resp.StatusCode)

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    assert.NotEmpty(t, result["user_id"])
}

// Тесты для каждого endpoint:
// POST /auth/register → 201, 409 (duplicate), 422 (validation)
// POST /auth/login    → 200, 401 (wrong pass), 403 (not verified), 429 (brute-force)
// POST /auth/refresh  → 200, 401 (expired), 401 (revoked)
// POST /auth/logout   → 200, 401 (invalid token)
```

#### 3.3 Чек-лист тестирования модуля

| Сценарий | Тип теста | Статус |
|----------|-----------|--------|
| Регистрация с валидными данными | Integration | ⬜ |
| Регистрация с занятым email | Integration | ⬜ |
| Регистрация со слабым паролем | Unit | ⬜ |
| Вход с верными данными | Integration | ⬜ |
| Вход с неверным паролем | Integration | ⬜ |
| Блокировка после 5 попыток | Integration | ⬜ |
| Обновление access token | Integration | ⬜ |
| Использование отозванного refresh | Integration | ⬜ |
| OAuth Google flow | E2E | ⬜ |
| Включение/выключение 2FA | Integration | ⬜ |
| Восстановление пароля | Integration | ⬜ |
| Нагрузочный тест на `/login` (500 rps) | Load | ⬜ |
| OWASP проверка SQL Injection | Security | ⬜ |

---

## Модуль 2: Профиль пользователя

### Фаза 1: Планирование задач

**Длительность:** 1 день  
**Участники:** Backend Developer, UI/UX Designer

#### 1.1 Задачи планирования

- [ ] Определить поля профиля и правила валидации
- [ ] Выбрать S3-совместимое хранилище для аватаров (Yandex Object Storage)
- [ ] Определить допустимые форматы/размеры аватара (JPG/PNG/WEBP, max 5MB)
- [ ] Спроектировать процедуру смены email (с переподтверждением)
- [ ] Определить требования GDPR/152-ФЗ для экспорта данных
- [ ] Согласовать формат экспорта: ZIP с JSON-файлами
- [ ] Создать тикеты USER-001..009

#### 1.2 Объём данных для GDPR-экспорта

```
export_user_12345.zip
├── profile.json          # Данные профиля
├── subscriptions.json    # Все подписки
├── payment_history.json  # История платежей
├── notifications.json    # История уведомлений
├── usage_tracking.json   # История использования
└── README.txt            # Описание формата
```

---

### Фаза 2: Разработка

#### 2.1 USER-003: Загрузка аватара

```go
// internal/users/service.go

func (s *UserService) UploadAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) error {
    // 1. Проверить MIME-тип (только image/jpeg, image/png, image/webp)
    // 2. Проверить размер (max 5MB)
    // 3. Ресайзить до 256x256 через imaging
    // 4. Загрузить в Object Storage: avatars/{userID}/{timestamp}.webp
    // 5. Удалить старый аватар если есть
    // 6. Обновить avatar_url в users
}
```

**Зависимости:**
- `github.com/disintegration/imaging` — ресайзинг
- `github.com/aws/aws-sdk-go-v2/service/s3` — Yandex Object Storage (S3-совместимый)

#### 2.2 USER-007: Экспорт данных

```go
// Генерация ZIP-архива со всеми данными пользователя
func (s *UserService) ExportData(ctx context.Context, userID string) (io.Reader, error) {
    // 1. Собрать данные из всех сервисов параллельно (errgroup)
    // 2. Создать ZIP в памяти
    // 3. Записать JSON файлы в ZIP
    // 4. Вернуть reader (отправить как attachment)
    // Операция может занять до 30 сек → ставить в очередь и отправлять на email
}
```

---

### Фаза 3: Тестирование

#### 3.1 Unit/Integration тесты

```go
func TestUpdateProfile_Success(t *testing.T) { ... }
func TestUpdateProfile_InvalidPhone(t *testing.T) { ... }
func TestUploadAvatar_ValidImage(t *testing.T) { ... }
func TestUploadAvatar_TooLarge(t *testing.T) { ... }
func TestUploadAvatar_InvalidMIME(t *testing.T) { ... }
func TestChangePassword_CorrectOldPassword(t *testing.T) { ... }
func TestChangePassword_WrongOldPassword(t *testing.T) { ... }
func TestChangeEmail_SendsVerification(t *testing.T) { ... }
func TestDeleteAccount_CascadeDeletes(t *testing.T) { ... }
func TestExportData_ContainsAllEntities(t *testing.T) { ... }
```

#### 3.2 Чек-лист тестирования модуля

| Сценарий | Тип | Статус |
|----------|-----|--------|
| Редактирование имени, фамилии | Integration | ⬜ |
| Загрузка PNG аватара | Integration | ⬜ |
| Загрузка файла > 5MB | Integration | ⬜ |
| Загрузка PDF вместо картинки | Integration | ⬜ |
| Смена пароля с верным старым | Integration | ⬜ |
| Смена email: отправка кода | Integration | ⬜ |
| Экспорт данных: ZIP содержит все файлы | Integration | ⬜ |
| Удаление аккаунта: каскадное удаление подписок | Integration | ⬜ |

---

## Модуль 3: Управление подписками

### Фаза 1: Планирование задач

**Длительность:** 2 дня  
**Участники:** Tech Lead, Backend Developer, Mobile Developer, Frontend Developer

#### 1.1 Задачи планирования

- [ ] Определить все поля сущности `subscription` и допустимые значения
- [ ] Согласовать список биллинговых циклов: `daily, weekly, monthly, quarterly, yearly, one_time`
- [ ] Определить список валют (RUB приоритет, USD, EUR, KZT, BYN)
- [ ] Согласовать формат дат (ISO 8601, UTC)
- [ ] Проработать состояния подписки (FSM): `active → paused → active`, `active → cancelled`
- [ ] Упорядочить сортировку по умолчанию: `next_billing_date ASC`
- [ ] Создать тикеты SUB-001..014
- [ ] Нарисовать wireframes (список, детали, форма добавления)

#### 1.2 Конечный автомат состояний

```
              ┌────────────────────────────────┐
              │               ACTIVE           │
              │  - Отображается в дашборде     │
              │  - Учитывается в аналитике     │
              │  - Генерирует уведомления      │
              └────┬─────────────────┬─────────┘
                   │ pause           │ cancel
                   ▼                 ▼
              ┌─────────┐       ┌──────────┐
              │ PAUSED  │       │CANCELLED │
              │         │       │          │
              └────┬────┘       └──────────┘
                   │ resume
                   ▼
              ┌─────────┐
              │ ACTIVE  │
              └─────────┘
```

#### 1.3 Категории подписок (seed data)

```json
[
  {"slug": "video", "name": "Видео стриминг", "icon": "🎬"},
  {"slug": "music", "name": "Музыка", "icon": "🎵"},
  {"slug": "software", "name": "Программное обеспечение", "icon": "💻"},
  {"slug": "cloud", "name": "Облачное хранилище", "icon": "☁️"},
  {"slug": "delivery", "name": "Доставка / Маркетплейсы", "icon": "📦"},
  {"slug": "games", "name": "Игры", "icon": "🎮"},
  {"slug": "news", "name": "СМИ / Чтение", "icon": "📰"},
  {"slug": "fitness", "name": "Спорт / Здоровье", "icon": "💪"},
  {"slug": "education", "name": "Образование", "icon": "📚"},
  {"slug": "other", "name": "Другое", "icon": "📌"}
]
```

---

### Фаза 2: Разработка

#### 2.1 SUB-001: Создание подписки

```go
// internal/subscriptions/service.go

type CreateSubscriptionRequest struct {
    Name            string   `json:"name"        validate:"required,min=1,max=255"`
    Price           float64  `json:"price"       validate:"required,gt=0"`
    Currency        string   `json:"currency"    validate:"required,oneof=RUB USD EUR KZT BYN"`
    BillingCycle    string   `json:"billing_cycle" validate:"required,oneof=daily weekly monthly quarterly yearly one_time"`
    NextBillingDate string   `json:"next_billing_date" validate:"required,datetime=2006-01-02"`
    StartDate       *string  `json:"start_date"`
    CategoryID      string   `json:"category_id" validate:"required,uuid"`
    ServiceID       *string  `json:"service_id"`
    Description     *string  `json:"description"`
    AutoRenew       bool     `json:"auto_renew"`
}

func (s *SubscriptionService) Create(ctx context.Context, userID string, req CreateSubscriptionRequest) (*Subscription, error) {
    // 1. Валидировать входные данные
    // 2. Проверить существование category_id
    // 3. Если service_id указан - подтянуть данные из каталога
    // 4. Создать запись в БД
    // 5. Запланировать уведомление о ближайшем списании
    // 6. Инвалидировать кэш аналитики пользователя
}
```

#### 2.2 SUB-006: Фильтрация и пагинация

```go
type ListSubscriptionsQuery struct {
    Status     []string `query:"status"`       // active, paused, cancelled
    CategoryID *string  `query:"category_id"`
    Search     *string  `query:"search"`       // fulltext по name
    SortBy     string   `query:"sort_by"`      // name, price, next_billing_date, created_at
    SortDir    string   `query:"sort_dir"`     // asc, desc
    Page       int      `query:"page"`
    PerPage    int      `query:"per_page"`
    Currency   *string  `query:"currency"`     // фильтр по валюте
}
// Формирование динамического SQL запроса через squirrel или pgx
```

#### 2.3 Расчёт годовой стоимости

```go
// Нормализация любого цикла к годовой стоимости
func AnnualCost(price float64, cycle string) float64 {
    multipliers := map[string]float64{
        "daily":     365,
        "weekly":    52,
        "monthly":   12,
        "quarterly": 4,
        "yearly":    1,
        "one_time":  0, // не учитывается в прогнозе
    }
    return price * multipliers[cycle]
}
```

---

### Фаза 3: Тестирование

#### 3.1 Unit-тесты

```go
func TestAnnualCost_Monthly(t *testing.T) {
    assert.Equal(t, 1200.0, AnnualCost(100, "monthly"))
}
func TestAnnualCost_Yearly(t *testing.T) {
    assert.Equal(t, 1200.0, AnnualCost(1200, "yearly"))
}
func TestCreateSubscription_Valid(t *testing.T) { ... }
func TestCreateSubscription_InvalidCategory(t *testing.T) { ... }
func TestPauseSubscription_FromActive(t *testing.T) { ... }
func TestPauseSubscription_AlreadyPaused(t *testing.T) { ... }
func TestResumeSubscription_FromPaused(t *testing.T) { ... }
func TestFilterSubscriptions_ByStatus(t *testing.T) { ... }
func TestFilterSubscriptions_ByCategory(t *testing.T) { ... }
func TestSearchSubscriptions_ByName(t *testing.T) { ... }
```

#### 3.2 Чек-лист тестирования

| Сценарий | Тип | Статус |
|----------|-----|--------|
| Создание подписки с валидными данными | Integration | ⬜ |
| Создание с несуществующей категорией | Integration | ⬜ |
| Список подписок с пагинацией | Integration | ⬜ |
| Фильтрация по статусу active | Integration | ⬜ |
| Сортировка по next_billing_date | Integration | ⬜ |
| Пауза активной подписки | Integration | ⬜ |
| Пауза уже паузированной → ошибка | Integration | ⬜ |
| Массовое удаление | Integration | ⬜ |
| Поиск по частичному названию | Integration | ⬜ |
| Чужая подписка → 403 | Integration | ⬜ |

---

## Модуль 4: Интеграция с банками

### Фаза 1: Планирование задач

**Длительность:** 3 дня  
**Участники:** Tech Lead, Backend Developer (senior), Security Officer

#### 1.1 Задачи планирования

- [ ] Изучить Open Banking API Сбербанка (OpenAPI spec, OAuth 2.0 flow)
- [ ] Изучить Tinkoff Open API (доступные endpoints для транзакций)
- [ ] Изучить ЮMoney API (payments history)
- [ ] Определить стратегию шифрования токенов: AES-256-GCM + AWS KMS / Yandex KMS
- [ ] Составить список MCC-кодов и паттернов мерчантов для распознавания подписок
- [ ] Определить частоту синхронизации: каждые 6 часов / при входе пользователя
- [ ] Согласовать UX для получения согласия пользователя (consent screen)
- [ ] Создать тикеты BANK-001..013
- [ ] Разработать план обработки ошибок API (401, 429, 5xx)

#### 1.2 Паттерны распознавания подписок в транзакциях

```go
// Паттерны мерчантов (регулярные выражения)
var subscriptionPatterns = []MerchantPattern{
    {Regex: `(?i)netflix`, ServiceName: "Netflix", Category: "video"},
    {Regex: `(?i)spotify`, ServiceName: "Spotify", Category: "music"},
    {Regex: `(?i)youtube\s*premium`, ServiceName: "YouTube Premium", Category: "video"},
    {Regex: `(?i)kinopoisk|кинопоиск`, ServiceName: "Кинопоиск", Category: "video"},
    {Regex: `(?i)yandex[\s\-]plus|яндекс[\s\-]плюс`, ServiceName: "Яндекс Плюс", Category: "delivery"},
    {Regex: `(?i)sber\s*prime|сбер\s*прайм`, ServiceName: "СберПрайм", Category: "delivery"},
    {Regex: `(?i)adobe`, ServiceName: "Adobe", Category: "software"},
    {Regex: `(?i)microsoft|office\s*365`, ServiceName: "Microsoft 365", Category: "software"},
    {Regex: `(?i)icloud`, ServiceName: "iCloud", Category: "cloud"},
    {Regex: `(?i)google\s*one`, ServiceName: "Google One", Category: "cloud"},
    // ... 100+ паттернов
}
```

---

### Фаза 2: Разработка

#### 2.1 BANK-010: Шифрование токенов (Critical)

```go
// pkg/crypto/token_vault.go

type TokenVault struct {
    kmsClient KMSClient
    keyID     string
}

func (v *TokenVault) Encrypt(plaintext string) (string, error) {
    // 1. Генерировать data encryption key (DEK) через KMS
    // 2. Зашифровать DEK мастер-ключом из KMS (envelope encryption)
    // 3. Зашифровать токен DEK (AES-256-GCM)
    // 4. Сохранить encrypted_dek + ciphertext + nonce вместе
    // Формат: base64(encrypted_dek_len + encrypted_dek + nonce + ciphertext)
}

func (v *TokenVault) Decrypt(ciphertext string) (string, error) {
    // 1. Распарсить ciphertext
    // 2. Расшифровать DEK через KMS
    // 3. Расшифровать токен DEK
}
```

#### 2.2 BANK-005/006: Алгоритм сопоставления транзакций

```
Алгоритм сопоставления (matching pipeline):

Для каждой транзакции:
1. Проверить список известных паттернов (regexp) → 80% случаев
2. Если не совпало → нормализовать merchant name (убрать спецсимволы, lower)
3. Проверить Elasticsearch индекс сервисов (fuzzy match, score > 0.8)
4. Если score 0.6-0.8 → отправить пользователю на подтверждение
5. Если < 0.6 → пометить как "unknown subscription"
6. Сохранить результат в subscription_transaction_matches с confidence_score
7. Если уверенность > 0.9 → автоматически привязать к подписке
```

#### 2.3 BANK-008: Периодическая синхронизация (Scheduler)

```go
// Использование: robfig/cron v3
func SetupBankSyncScheduler(svc BankService) {
    c := cron.New()
    // Каждые 6 часов синхронизировать все активные подключения
    c.AddFunc("0 */6 * * *", func() {
        connections := svc.GetAllActiveConnections(ctx)
        // Обработка в воркер-пуле (max 10 параллельных)
        pool := workerpool.New(10)
        for _, conn := range connections {
            conn := conn
            pool.Submit(func() {
                svc.SyncConnection(ctx, conn.ID)
            })
        }
        pool.StopWait()
    })
    c.Start()
}
```

---

### Фаза 3: Тестирование

#### 3.1 Тестирование с mock-банком

```go
// Мок сервера банка для тестов
func SetupMockBankServer() *httptest.Server {
    mux := http.NewServeMux()
    mux.HandleFunc("/oauth/token", handleMockToken)
    mux.HandleFunc("/v1/transactions", handleMockTransactions)
    return httptest.NewServer(mux)
}

func TestSberbankConnect_Success(t *testing.T) { ... }
func TestTinkoffSync_ParsesTransactions(t *testing.T) { ... }
func TestMatchTransaction_Netflix(t *testing.T) { ... }
func TestMatchTransaction_Unknown(t *testing.T) { ... }
func TestTokenEncryption_Decrypt(t *testing.T) { ... }
func TestTokenExpiry_Refresh(t *testing.T) { ... }
func TestBankSync_RateLimitHandling(t *testing.T) { ... }
func TestBankSync_NetworkTimeout(t *testing.T) { ... }
```

#### 3.2 Чек-лист тестирования

| Сценарий | Тип | Статус |
|----------|-----|--------|
| OAuth flow Сбербанк (mock) | Integration | ⬜ |
| OAuth flow Тинькофф (mock) | Integration | ⬜ |
| Парсинг транзакции Netflix | Unit | ⬜ |
| Парсинг неизвестной транзакции | Unit | ⬜ |
| Шифрование/расшифровка токена | Unit | ⬜ |
| Ротация истёкшего токена | Integration | ⬜ |
| Обработка 429 от банка (retry + backoff) | Unit | ⬜ |
| Отключение банка очищает токены | Integration | ⬜ |
| Синхронизация не создаёт дубли | Integration | ⬜ |
| Security: токен не логируется в plaintext | Security | ⬜ |

---

## Модуль 5: Парсинг почты

### Фаза 1: Планирование задач

**Длительность:** 3 дня  
**Участники:** Tech Lead, Backend Developer, ML Engineer

#### 1.1 Задачи планирования

- [ ] Изучить Gmail API (users.messages.list, users.messages.get)
- [ ] Изучить Yandex Mail OAuth + IMAP
- [ ] Изучить Mail.ru OAuth + IMAP
- [ ] Сформировать список email-отправителей сервисов (500+)
- [ ] Определить поисковые фильтры для каждого провайдера
- [ ] Согласовать UX consent: показывать какие именно данные будут считаны
- [ ] Выбрать NLP-подход: template matching + ML classifier
- [ ] Собрать датасет примеров чеков (минимум 50 примеров на сервис)
- [ ] Создать тикеты EMAIL-001..016

#### 1.2 Стратегия парсинга писем

```
Уровень 1: Sender-based фильтрация
  • Проверить отправителя по whitelist (no-reply@netflix.com, и т.д.)
  
Уровень 2: Subject-based фильтрация  
  • Регулярки: "чек", "квитанция", "receipt", "invoice", "списание", "платёж"
  
Уровень 3: HTML-парсинг структуры
  • Найти блок с суммой: регулярки для "₽ 1,490", "1490 RUB", "$9.99"
  • Найти дату: ISO даты, русский формат "14 февраля 2026"
  • Найти название сервиса: H1/H2 заголовки, мета-теги
  
Уровень 4: ML-классификатор (fallback)
  • Модель: distilbert-multilingual (RU + EN)
  • Fine-tuned на датасете чеков
  • Confidence threshold: 0.75
```

---

### Фаза 2: Разработка

#### 2.1 EMAIL-007/008/009: Парсер HTML писем

```go
// internal/emailparser/parser.go

type ParsedReceipt struct {
    ServiceName string
    Amount      float64
    Currency    string
    Date        time.Time
    Source      string    // email subject + sender
    Confidence  float64   // 0.0 - 1.0
    RawHTML     string
}

type EmailParser struct {
    patterns    []ReceiptPattern
    mlClient    MLClient
}

func (p *EmailParser) Parse(html, sender, subject string) (*ParsedReceipt, error) {
    // 1. Попытка 1: паттерн по отправителю
    if pattern := p.findPatternBySender(sender); pattern != nil {
        receipt, err := pattern.Extract(html)
        if err == nil && receipt.Confidence > 0.9 {
            return receipt, nil
        }
    }

    // 2. Попытка 2: общий RegExp парсер
    receipt := p.genericParse(html, subject)
    if receipt.Confidence > 0.75 {
        return receipt, nil
    }

    // 3. Попытка 3: ML-классификатор
    return p.mlClient.Classify(html, subject)
}

// Пример шаблона для Netflix
var NetflixPattern = ReceiptPattern{
    Sender:  "info@account.netflix.com",
    Amount:  regexp.MustCompile(`(\d[\d\s,.]+)\s*(руб|RUB|₽|\$|USD|€|EUR)`),
    Date:    regexp.MustCompile(`(\d{1,2})\s+(января|февраля|марта|апреля|мая|июня|июля|августа|сентября|октября|ноября|декабря)\s+(\d{4})`),
    Service: "Netflix",
}
```

#### 2.2 EMAIL-013: Периодическая проверка новых писем

```go
// Инкрементальная синхронизация через historyId (Gmail) / UID (IMAP)
func (s *EmailService) IncrementalSync(ctx context.Context, connID string) error {
    conn := s.repo.GetConnection(ctx, connID)
    
    var newMessages []Message
    switch conn.Provider {
    case "gmail":
        // Использовать history.list с historyId с последней синхронизации
        newMessages, _ = s.gmailClient.GetMessagesSince(conn.LastHistoryID)
        conn.LastHistoryID = latestHistoryID
    case "yandex", "mailru":
        // IMAP FETCH с UID > lastUID
        newMessages, _ = s.imapClient.FetchSince(conn.LastUID)
    }
    
    for _, msg := range newMessages {
        receipt, err := s.parser.Parse(msg.HTML, msg.Sender, msg.Subject)
        if err != nil || receipt.Confidence < 0.5 { continue }
        s.repo.SaveParsedReceipt(ctx, conn.UserID, receipt)
    }
    conn.LastSyncAt = time.Now()
    return s.repo.UpdateConnection(ctx, conn)
}
```

---

### Фаза 3: Тестирование

#### 3.1 Тесты парсера на реальных примерах

```go
// testdata/ — папка с реальными HTML письмами (анонимизированными)
func TestNetflixReceiptParsing(t *testing.T) {
    html := readTestFile(t, "testdata/netflix_receipt_ru.html")
    parser := NewEmailParser(patterns, nil)
    
    receipt, err := parser.Parse(html, "info@account.netflix.com", "Ваш чек Netflix")
    
    require.NoError(t, err)
    assert.Equal(t, "Netflix", receipt.ServiceName)
    assert.Equal(t, 1490.0, receipt.Amount)
    assert.Equal(t, "RUB", receipt.Currency)
    assert.True(t, receipt.Confidence > 0.9)
}

// Тесты для каждого сервиса в testdata/:
// netflix, spotify, youtube_premium, kinopoisk, okko, ivi,
// adobe, microsoft365, jetbrains, yandex_disk, icloud,
// google_one, yandex_plus, sber_prime, wildberries_premium...
```

#### 3.2 Чек-лист тестирования

| Сценарий | Тип | Статус |
|----------|-----|--------|
| Парсинг чека Netflix | Unit | ⬜ |
| Парсинг чека Spotify | Unit | ⬜ |
| Парсинг чека Adobe (USD) | Unit | ⬜ |
| Парсинг с низкой уверенностью → на ревью | Unit | ⬜ |
| Gmail OAuth connect + fetch (mock API) | Integration | ⬜ |
| IMAP Яндекс.Почта (mock server) | Integration | ⬜ |
| Инкрементальная синхронизация (только новые) | Integration | ⬜ |
| Дублирование чека не создаёт вторую подписку | Integration | ⬜ |
| Пользователь отклоняет распознанную подписку | Integration | ⬜ |
| Согласие отозвано → токены удаляются | Integration | ⬜ |
| Парсер не читает личную переписку | Security | ⬜ |

---

## Модуль 6: Аналитика

### Фаза 1: Планирование задач

**Длительность:** 2 дня  
**Участники:** Backend Developer, Frontend Developer, UI/UX Designer

#### 1.1 Задачи планирования

- [ ] Определить набор метрик для дашборда (KPIs)
- [ ] Спроектировать схему материализованных представлений в ClickHouse
- [ ] Определить стратегию кэширования (TTL по типу запроса)
- [ ] Нарисовать wireframes всех графиков
- [ ] Согласовать форматы дат в API ответах
- [ ] Определить политику пересчёта (при изменении подписки → инвалидировать кэш)
- [ ] Создать тикеты ANAL-001..012

#### 1.2 Кэш-стратегия

```
Тип запроса               | TTL   | Инвалидация
─────────────────────────────────────────────────
/analytics/overview        | 1 ч   | При изменении подписки
/analytics/spending/monthly| 1 ч   | При изменении подписки
/analytics/spending/categories | 1 ч | При изменении подписки
/analytics/top-subscriptions | 1 ч  | При изменении подписки
/analytics/comparison       | 6 ч  | Нет (историческая)
```

---

### Фаза 2: Разработка

#### 2.1 ANAL-001: Overview endpoint

```go
type AnalyticsOverview struct {
    TotalMonthly     float64 `json:"total_monthly"`     // Сумма за текущий месяц
    TotalYearly      float64 `json:"total_yearly"`      // Прогноз на год
    ActiveCount      int     `json:"active_count"`      // Кол-во активных
    MostExpensive    *Sub    `json:"most_expensive"`    // Самая дорогая
    NextBillingIn    int     `json:"next_billing_in"`   // Дней до ближайшего списания
    NextBillingAmount float64 `json:"next_billing_amount"` // Сумма ближайшего
    ChangePercent    float64 `json:"change_percent"`    // Изменение vs прошлый месяц
    Currency         string  `json:"currency"`          // Основная валюта
}

func (s *AnalyticsService) GetOverview(ctx context.Context, userID string) (*AnalyticsOverview, error) {
    cacheKey := fmt.Sprintf("analytics:overview:%s", userID)
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        return cached, nil
    }
    
    result := s.repo.ComputeOverview(ctx, userID)
    s.cache.Set(ctx, cacheKey, result, time.Hour)
    return result, nil
}
```

#### 2.2 ANAL-008: Генерация PDF отчёта

```go
// Использование: johnfercher/maroto (PDF для Go)
func (s *AnalyticsService) GeneratePDFReport(ctx context.Context, userID string, period Period) ([]byte, error) {
    data := s.collectReportData(ctx, userID, period)
    
    m := maroto.New(maroto.Portrait, maroto.A4)
    // Добавить заголовок, таблицы, диаграммы
    // Диаграммы генерировать через go-chart → PNG → вставить в PDF
    
    return m.Output()
}
```

---

### Фаза 3: Тестирование

```go
func TestOverview_CorrectTotals(t *testing.T) {
    // Создать 3 подписки: 500 RUB/month, 1000 RUB/month, 12000 RUB/year
    // Ожидаемый total_monthly = 500 + 1000 + 1000 = 2500
    // Ожидаемый total_yearly = 6000 + 12000 + 12000 = 30000
}
func TestSpendingByCategory_Grouping(t *testing.T) { ... }
func TestComparison_MoreThanPrevMonth(t *testing.T) { ... }
func TestComparison_NoPreviousData(t *testing.T) { ... }
func TestCacheInvalidation_OnSubscriptionUpdate(t *testing.T) { ... }
func TestExportCSV_ContainsAllSubscriptions(t *testing.T) { ... }
func TestExportPDF_ValidDocument(t *testing.T) { ... }
```

#### 3.2 Чек-лист тестирования

| Сценарий | Тип | Статус |
|----------|-----|--------|
| Overview: правильный подсчёт total | Unit | ⬜ |
| Overview: конвертация валют в RUB | Unit | ⬜ |
| График по месяцам: 12 точек | Integration | ⬜ |
| Фильтр по периоду (последние 3 мес.) | Integration | ⬜ |
| Кэш возвращает данные без SQL | Unit | ⬜ |
| Инвалидация кэша при добавлении подписки | Integration | ⬜ |
| CSV экспорт: заголовки + данные | Integration | ⬜ |
| PDF генерация без ошибок | Integration | ⬜ |
| Нагрузочный тест: 200 rps на /overview | Load | ⬜ |

---

## Модуль 7: Отслеживание использования

### Фаза 1: Планирование задач

**Длительность:** 1 день  
**Участники:** Backend Developer, Mobile Developer

#### 1.1 Задачи планирования

- [ ] Определить типы использования: `opened`, `watched`, `listened`, `read`, `used`
- [ ] Определить метрику "неактивности": N дней без отметки
- [ ] Согласовать расчёт "стоимости за использование": `monthly_cost / usages_this_month`
- [ ] Определить UI quick-action: свайп или кнопка на карточке подписки
- [ ] Создать тикеты USE-001..010

---

### Фаза 2: Разработка

#### 2.1 USE-005: Расчёт стоимости за использование

```go
type UsageStats struct {
    TotalUsages      int     `json:"total_usages"`
    UsagesThisMonth  int     `json:"usages_this_month"`
    LastUsedAt       *time.Time `json:"last_used_at"`
    DaysSinceLastUse *int    `json:"days_since_last_use"`
    CostPerUse       float64 `json:"cost_per_use"`      // monthly_price / usages_this_month
    IsUnderutilized  bool    `json:"is_underutilized"`  // < 2 раза в месяц
}

func (s *UsageService) GetStats(ctx context.Context, subID, userID string) (*UsageStats, error) {
    sub := s.subRepo.Get(ctx, subID)
    usages := s.repo.GetThisMonth(ctx, subID)
    
    costPerUse := 0.0
    if len(usages) > 0 {
        costPerUse = sub.MonthlyPrice() / float64(len(usages))
    }
    
    return &UsageStats{
        UsagesThisMonth: len(usages),
        CostPerUse:      costPerUse,
        IsUnderutilized: len(usages) < 2,
    }, nil
}
```

---

### Фаза 3: Тестирование

```go
func TestMarkUsage_CreatesRecord(t *testing.T) { ... }
func TestMarkUsage_TwiceInOneDay_Allowed(t *testing.T) { ... }
func TestCostPerUse_NoUsages_ReturnsZero(t *testing.T) { ... }
func TestCostPerUse_5Usages_Calculates(t *testing.T) { ... }
func TestInactivityDetection_30Days(t *testing.T) { ... }
func TestCalendar_ShowsUsageDots(t *testing.T) { ... }
```

| Сценарий | Тип | Статус |
|----------|-----|--------|
| Отметка использования создаёт запись | Integration | ⬜ |
| Дважды в один день — допустимо | Integration | ⬜ |
| CostPerUse при 4 отметках в месяц | Unit | ⬜ |
| Подписка не использовалась 30+ дней | Unit | ⬜ |
| Удаление отметки | Integration | ⬜ |
| Чужая подписка → 403 | Integration | ⬜ |

---

## Модуль 8: Прогнозирование

### Фаза 1: Планирование задач

**Длительность:** 2 дня  
**Участники:** Tech Lead, ML Engineer, Backend Developer

#### 1.1 Задачи планирования

- [ ] Определить алгоритм базового прогноза (детерминированный)
- [ ] Определить источники данных для тренда цен (история payment_history)
- [ ] Согласовать формат ответа: помесячный массив + итоговая сумма
- [ ] Спроектировать симулятор сценариев "что если"
- [ ] Определить логику бюджетных алертов
- [ ] Создать тикеты PRED-001..010

#### 1.2 Алгоритм прогноза (v1 — детерминированный)

```
ДЛЯ каждого месяца в [текущий_месяц .. текущий_месяц + 12]:
  total = 0
  ДЛЯ каждой активной подписки:
    ЕСЛИ подписка будет списана в этом месяце:
      price = base_price * price_change_coefficient
      total += convert_to_rub(price, currency)
  monthly_forecast[month] = total

price_change_coefficient:
  - Если есть >= 3 записей в payment_history: regression slope
  - Иначе: 1.0 (нет изменения)
```

---

### Фаза 2: Разработка

#### 2.1 PRED-005: Симулятор сценариев

```go
type SimulationRequest struct {
    AddSubscriptions    []SimSub `json:"add_subscriptions"`
    RemoveSubscriptions []string `json:"remove_subscriptions"` // IDs
    ChangePrice         []PriceChange `json:"change_prices"`
}

type SimulationResponse struct {
    BaseMonthly    []MonthForecast `json:"base_monthly"`    // Без изменений
    SimMonthly     []MonthForecast `json:"sim_monthly"`     // С изменениями
    BasedYearly    float64         `json:"base_yearly"`
    SimYearly      float64         `json:"sim_yearly"`
    Difference     float64         `json:"difference"`      // Разница в год
}
```

---

### Фаза 3: Тестирование

```go
func TestYearlyForecast_MonthlySubscription(t *testing.T) {
    // 1 подписка 500 RUB/month → 6000 RUB/year
}
func TestYearlyForecast_YearlySubscription(t *testing.T) {
    // 1 подписка 1200 RUB/year, следующий платёж через 2 мес
    // → прогноз: 0*10 + 1200 один раз = 1200 за год
}
func TestForecast_MixedCurrencies_ConvertedToRUB(t *testing.T) { ... }
func TestForecast_PausedSubscriptionExcluded(t *testing.T) { ... }
func TestSimulation_AddSubscription_IncreasesTotal(t *testing.T) { ... }
func TestSimulation_RemoveSubscription_DecreasesTotal(t *testing.T) { ... }
func TestBudgetAlert_TriggeredWhenExceeded(t *testing.T) { ... }
```

| Сценарий | Тип | Статус |
|----------|-----|--------|
| Прогноз одной monthly подписки | Unit | ⬜ |
| Прогноз yearly подписки | Unit | ⬜ |
| Паузированные подписки не учитываются | Unit | ⬜ |
| Конвертация USD → RUB по курсу | Unit | ⬜ |
| Симулятор: добавление подписки | Integration | ⬜ |
| Бюджетный алерт: превышение лимита | Integration | ⬜ |

---

## Модуль 9: Уведомления

### Фаза 1: Планирование задач

**Длительность:** 2 дня  
**Участники:** Backend Developer, Mobile Developer, Frontend Developer

#### 1.1 Задачи планирования

- [ ] Составить полный список типов уведомлений с текстами
- [ ] Выбрать push-провайдер: FCM для Android, APNs через FCM для iOS
- [ ] Определить приоритеты доставки: billing_reminder=HIGH, inactivity=NORMAL
- [ ] Согласовать логику тихих часов
- [ ] Спроектировать таблицу очереди уведомлений
- [ ] Определить retry-политику: 3 попытки с exponential backoff
- [ ] Создать тикеты NOTIF-001..015

#### 1.2 Очередь уведомлений

```sql
CREATE TABLE notification_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    channel VARCHAR(20) NOT NULL,  -- push, email, in_app
    payload JSONB NOT NULL,
    scheduled_at TIMESTAMP NOT NULL,
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    last_attempt_at TIMESTAMP,
    error_message TEXT,
    status VARCHAR(20) DEFAULT 'pending',  -- pending, processing, sent, failed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_queue_scheduled (status, scheduled_at) WHERE status = 'pending'
);
```

---

### Фаза 2: Разработка

#### 2.1 NOTIF-004: Триггер предстоящего списания

```go
// Джоб запускается каждый час
func (s *NotificationService) ScheduleBillingReminders(ctx context.Context) error {
    // Найти подписки, где next_billing_date IN (now+1d, now+3d, now+7d)
    // и для которых ещё не создано уведомление за этот период
    
    daysToNotify := []int{1, 3, 7}
    for _, days := range daysToNotify {
        targetDate := time.Now().AddDate(0, 0, days).Truncate(24 * time.Hour)
        subs := s.subRepo.FindByNextBillingDate(ctx, targetDate)
        
        for _, sub := range subs {
            settings := s.settingsRepo.Get(ctx, sub.UserID)
            if !containsInt(settings.BillingReminderDays, days) { continue }
            
            s.enqueue(ctx, &NotificationJob{
                UserID: sub.UserID,
                Type:   "billing_reminder",
                Payload: map[string]any{
                    "subscription_name": sub.Name,
                    "amount":           sub.Price,
                    "currency":         sub.Currency,
                    "days":             days,
                    "billing_date":     sub.NextBillingDate,
                },
                ScheduledAt: time.Now(),
            })
        }
    }
    return nil
}
```

#### 2.2 NOTIF-001: FCM Worker

```go
// workers/notification_worker.go
func (w *NotificationWorker) ProcessQueue(ctx context.Context) {
    for {
        jobs := w.queue.Dequeue(ctx, batchSize=100)
        for _, job := range jobs {
            var err error
            switch job.Channel {
            case "push":
                err = w.fcm.Send(ctx, job.DeviceToken, job.Payload)
            case "email":
                err = w.emailSvc.Send(ctx, job.UserEmail, job.Payload)
            case "in_app":
                err = w.inAppSvc.Save(ctx, job.UserID, job.Payload)
            }
            
            if err != nil {
                job.Attempts++
                if job.Attempts >= job.MaxAttempts {
                    w.queue.MarkFailed(ctx, job.ID, err.Error())
                } else {
                    // Exponential backoff: 1m, 5m, 25m
                    delay := time.Duration(math.Pow(5, float64(job.Attempts))) * time.Minute
                    w.queue.Reschedule(ctx, job.ID, time.Now().Add(delay))
                }
            } else {
                w.queue.MarkSent(ctx, job.ID)
            }
        }
        time.Sleep(5 * time.Second)
    }
}
```

---

### Фаза 3: Тестирование

```go
func TestBillingReminder_ScheduledForCorrectDays(t *testing.T) { ... }
func TestBillingReminder_RespectUserSettings(t *testing.T) { ... }
func TestQuietHours_NotSentDuringNight(t *testing.T) { ... }
func TestRetry_FailedPushRetried3Times(t *testing.T) { ... }
func TestRetry_AfterMaxAttempts_MarkedFailed(t *testing.T) { ... }
func TestInactivityTrigger_After30Days(t *testing.T) { ... }
func TestPriceChangeTrigger_WhenPriceUpdated(t *testing.T) { ... }
func TestMarkAsRead_UpdatesStatus(t *testing.T) { ... }
```

| Сценарий | Тип | Статус |
|----------|-----|--------|
| Напоминание создаётся за 3 дня | Integration | ⬜ |
| Напоминание за 7 дней | Integration | ⬜ |
| Не создаётся дубль при повторном запуске | Integration | ⬜ |
| Тихие часы: уведомление откладывается | Unit | ⬜ |
| Retry при ошибке FCM | Unit | ⬜ |
| Уведомление об изменении цены | Integration | ⬜ |
| Отключение push → только email | Integration | ⬜ |
| Нагрузочный тест: 10K уведомлений/сек | Load | ⬜ |

---

## Модуль 10: Каталог сервисов и альтернативы

### Фаза 1: Планирование задач

**Длительность:** 2 дня  
**Участники:** Product Manager, Backend Developer, Content Manager

#### 1.1 Задачи планирования

- [ ] Составить список первичных 500 сервисов с данными
- [ ] Определить структуру данных для ценовых планов (массив)
- [ ] Разработать алгоритм подбора альтернатив (по категории + диапазону цен)
- [ ] Определить источники данных о ценах (ручное + парсинг сайтов)
- [ ] Создать инструмент для наполнения каталога (admin CRUD)
- [ ] Создать тикеты CAT-001..005, ALT-001..006

#### 1.2 Алгоритм альтернатив

```
Входные данные: подписка пользователя (сервис X, цена P, категория C)

1. Найти все сервисы в категории C, кроме X
2. Фильтровать: есть план <= P * 1.2 (не дороже на 20%)
3. Рассчитать similarity_score:
   - Совпадение функций: +0.4
   - Близость цены: +0.3 * (1 - |price_diff| / P)
   - Рейтинг сервиса: +0.2 * (rating / 5)
   - Популярность у пользователей MultiSub: +0.1
4. Сортировать по similarity_score DESC
5. Вернуть топ-5
```

---

### Фаза 2: Разработка

#### 2.1 CAT-001: Seed данных каталога

```go
// migrations/seed_services.go
// Скрипт для начального наполнения 500+ сервисов
// Данные хранятся в JSON-файлах: data/services/*.json

type ServiceSeed struct {
    Name        string          `json:"name"`
    Category    string          `json:"category_slug"`
    LogoURL     string          `json:"logo_url"`
    WebsiteURL  string          `json:"website_url"`
    CancelURL   string          `json:"cancel_url"`
    Plans       []PricingPlan   `json:"plans"`
    Features    []string        `json:"features"`
    Country     []string        `json:"countries"`  // ["RU","US","EU"]
}
```

#### 2.2 ALT-001: Endpoint альтернатив

```go
func (h *CatalogHandler) GetAlternatives(c *fiber.Ctx) error {
    serviceID := c.Params("id")
    userSub := h.subRepo.FindByService(c.Context(), userID, serviceID)
    
    alternatives := h.altSvc.FindAlternatives(c.Context(), AlternativesQuery{
        ServiceID:   serviceID,
        UserBudget:  userSub.Price * 1.2,
        CategoryID:  userSub.CategoryID,
        Limit:       5,
    })
    
    // Обогатить данными о потенциальной экономии
    for i := range alternatives {
        alternatives[i].PotentialSavings = userSub.AnnualCost() - alternatives[i].AnnualCost()
    }
    
    return c.JSON(alternatives)
}
```

---

### Фаза 3: Тестирование

```go
func TestAlternatives_ReturnsSameCategory(t *testing.T) { ... }
func TestAlternatives_NoCheaperThan20Percent(t *testing.T) { ... }
func TestAlternatives_ExcludesCurrentService(t *testing.T) { ... }
func TestCatalogSearch_ByName(t *testing.T) { ... }
func TestCatalogSearch_CaseInsensitive(t *testing.T) { ... }
func TestSavingsCalculation_Correct(t *testing.T) { ... }
```

| Сценарий | Тип | Статус |
|----------|-----|--------|
| Альтернативы той же категории | Unit | ⬜ |
| Исключение текущего сервиса | Unit | ⬜ |
| Поиск без результатов — пустой массив | Integration | ⬜ |
| Расчёт экономии за год | Unit | ⬜ |
| Полнотекстовый поиск (Elasticsearch) | Integration | ⬜ |

---

## Модуль 11: Отмена подписок

### Фаза 1: Планирование задач

**Длительность:** 2 дня  
**Участники:** Product Manager, Content Manager, Backend Developer

#### 1.1 Задачи планирования

- [ ] Для каждого сервиса из каталога написать инструкцию по отмене (шаги)
- [ ] Собрать список direct cancel URLs для каждого сервиса
- [ ] Разработать шаблоны писем для поддержки (RU + EN)
- [ ] Определить список причин отмены (для аналитики)
- [ ] Создать тикеты CANC-001..010

#### 1.2 Структура инструкции по отмене

```json
{
  "service_id": "uuid-netflix",
  "cancel_methods": [
    {
      "method": "web",
      "title": "Через сайт",
      "url": "https://www.netflix.com/cancelplan",
      "steps": [
        "Войдите в аккаунт на netflix.com",
        "Нажмите на иконку профиля → «Аккаунт»",
        "В разделе «Членство и выставление счетов» нажмите «Отключить подписку»",
        "Подтвердите отмену"
      ],
      "difficulty": "easy",
      "time_minutes": 3
    },
    {
      "method": "email",
      "title": "Письмо в поддержку",
      "support_email": "support@netflix.com",
      "template_id": "cancel_subscription_ru"
    }
  ],
  "important_notes": [
    "Доступ сохраняется до конца оплаченного периода",
    "Скачанный контент удалится через 7 дней"
  ],
  "refund_policy": "Netflix не возвращает средства за текущий период"
}
```

---

### Фаза 2: Разработка

#### 2.1 CANC-003: Генератор письма

```go
type CancelLetterRequest struct {
    SubscriptionID string `json:"subscription_id"`
    Reason         string `json:"reason"`          // price, unused, alternative, other
    ReasonDetail   string `json:"reason_detail"`   // пользовательский текст
    Language       string `json:"language"`        // ru, en
}

func (s *CancellationService) GenerateCancelLetter(req CancelLetterRequest) *CancelLetter {
    tmpl := s.templates.Get(req.Language)
    user := s.userRepo.Get(req.SubscriptionID)
    sub := s.subRepo.Get(req.SubscriptionID)
    service := s.catalogRepo.GetBySubID(req.SubscriptionID)
    
    return &CancelLetter{
        To:      service.SupportEmail,
        Subject: fmt.Sprintf("Запрос на отмену подписки - %s", user.FullName()),
        Body:    tmpl.Render(user, sub, service, req.Reason),
    }
}
```

---

### Фаза 3: Тестирование

```go
func TestCancelGuide_HasSteps(t *testing.T) { ... }
func TestLetterGeneration_Russian(t *testing.T) { ... }
func TestLetterGeneration_English(t *testing.T) { ... }
func TestLetterContainsUserEmail(t *testing.T) { ... }
func TestMarkCancelled_ChangesStatus(t *testing.T) { ... }
func TestCancelURL_NotEmpty(t *testing.T) { ... }
```

---

## Модуль 12: Административная панель

### Фаза 1: Планирование задач

**Длительность:** 2 дня  
**Участники:** Tech Lead, Backend Developer

#### 1.1 Задачи планирования

- [ ] Определить роли: `super_admin`, `admin`, `moderator`, `support`
- [ ] Описать матрицу прав доступа (RBAC)
- [ ] Определить метрики дашборда платформы
- [ ] Спроектировать аудит-лог (immutable)
- [ ] Создать тикеты ADM-001..012

#### 1.2 Матрица прав (RBAC)

| Действие | super_admin | admin | moderator | support |
|----------|:-----------:|:-----:|:---------:|:-------:|
| Просмотр пользователей | ✅ | ✅ | ✅ | ✅ |
| Удаление пользователей | ✅ | ✅ | ❌ | ❌ |
| Управление каталогом | ✅ | ✅ | ✅ | ❌ |
| Просмотр аудит-лога | ✅ | ✅ | ❌ | ❌ |
| Настройки системы | ✅ | ❌ | ❌ | ❌ |
| Управление ролями | ✅ | ❌ | ❌ | ❌ |

---

### Фаза 2: Разработка

#### 2.1 ADM-007: Аудит-лог

```go
// Иммутабельный лог всех административных действий
type AuditEntry struct {
    ID         UUID      
    AdminID    UUID      
    Action     string    // "user.delete", "catalog.service.update", ...
    ResourceID string    // ID затронутого ресурса
    OldValue   JSONB     // Состояние до
    NewValue   JSONB     // Состояние после
    IPAddress  string    
    UserAgent  string    
    CreatedAt  time.Time
}

// Middleware для автоматического логирования
func AuditMiddleware(svc AuditService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        err := c.Next()
        if c.Method() != "GET" && err == nil {
            svc.Log(c.Context(), AuditEntry{...})
        }
        return err
    }
}
```

---

### Фаза 3: Тестирование

```go
func TestRBAC_ModeratorCannotDeleteUser(t *testing.T) { ... }
func TestRBAC_AdminCanManageCatalog(t *testing.T) { ... }
func TestAuditLog_CapturesDeleteAction(t *testing.T) { ... }
func TestAuditLog_Immutable_CannotUpdate(t *testing.T) { ... }
func TestPlatformStats_CountsActiveUsers(t *testing.T) { ... }
```

---

## Модуль 13: Мобильное приложение

### Фаза 1: Планирование задач

**Длительность:** 3 дня  
**Участники:** Mobile Developer ×2, UI/UX Designer, Tech Lead

#### 1.1 Задачи планирования

- [ ] Финализировать дизайн-систему в Figma (компоненты, токены)
- [ ] Выбрать Expo workflow: `bare workflow` (полный доступ к нативным модулям)
- [ ] Определить структуру навигации: Tab (4 вкладки) + Stack
- [ ] Проработать offline-стратегию: что кэшируется локально
- [ ] Спроектировать структуру Zustand stores
- [ ] Создать тикеты MOB-001..021

#### 1.2 Структура навигации

```
RootNavigator (Stack)
├── AuthNavigator (Stack) — unauthorised
│   ├── WelcomeScreen
│   ├── LoginScreen
│   ├── RegisterScreen
│   ├── ForgotPasswordScreen
│   └── OnboardingNavigator (Stack)
│       ├── Onboarding1Screen (connect bank)
│       ├── Onboarding2Screen (connect email)
│       └── Onboarding3Screen (add first sub)
└── MainNavigator (Bottom Tabs) — authorised
    ├── HomeTab → DashboardScreen
    ├── SubscriptionsTab → SubscriptionsListScreen
    │   ├── SubscriptionDetailScreen (Stack)
    │   ├── AddSubscriptionScreen (Stack/Modal)
    │   └── EditSubscriptionScreen (Stack/Modal)
    ├── AnalyticsTab → AnalyticsScreen
    └── ProfileTab → ProfileScreen
        ├── NotificationsScreen
        ├── IntegrationsScreen
        │   ├── BankConnectScreen
        │   └── EmailConnectScreen
        ├── SettingsScreen
        └── AboutScreen
```

#### 1.3 Структура Zustand

```typescript
// stores/authStore.ts
interface AuthStore {
  user: User | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshToken: () => Promise<void>;
}

// stores/subscriptionsStore.ts
interface SubscriptionsStore {
  subscriptions: Subscription[];
  isLoading: boolean;
  filters: FilterState;
  fetchSubscriptions: () => Promise<void>;
  createSubscription: (data: CreateSubDTO) => Promise<void>;
  updateSubscription: (id: string, data: UpdateSubDTO) => Promise<void>;
  deleteSubscription: (id: string) => Promise<void>;
  setFilters: (filters: Partial<FilterState>) => void;
}

// TanStack Query для серверного состояния
// Zustand только для клиентского UI-состояния
```

---

### Фаза 2: Разработка

#### 2.1 MOB-009: Экран аналитики (Victory Native)

```tsx
// screens/AnalyticsScreen.tsx
import { VictoryBar, VictoryPie, VictoryChart, VictoryAxis } from 'victory-native';

export const AnalyticsScreen = () => {
  const { data: overview } = useQuery(['analytics', 'overview'], fetchOverview);
  const { data: monthly } = useQuery(['analytics', 'monthly'], fetchMonthly);

  return (
    <ScrollView>
      <OverviewCards data={overview} />
      
      {/* График расходов по месяцам */}
      <VictoryChart>
        <VictoryAxis />
        <VictoryBar data={monthly?.data} x="month" y="total" />
      </VictoryChart>
      
      {/* Pie chart по категориям */}
      <VictoryPie data={overview?.categories} x="name" y="amount" />
      
      <UpcomingBillingList />
    </ScrollView>
  );
};
```

#### 2.2 MOB-016: Биометрическая аутентификация

```tsx
// hooks/useBiometrics.ts
import ReactNativeBiometrics from 'react-native-biometrics';

export const useBiometrics = () => {
  const authenticate = async (): Promise<boolean> => {
    const rnBiometrics = new ReactNativeBiometrics();
    const { available, biometryType } = await rnBiometrics.isSensorAvailable();
    
    if (!available) return false;
    
    const { success } = await rnBiometrics.simplePrompt({
      promptMessage: 'Подтвердите вход',
      cancelButtonText: 'Использовать пароль',
    });
    
    return success;
  };
  
  return { authenticate };
};
```

---

### Фаза 3: Тестирование

#### 3.1 Unit-тесты компонентов (Jest + React Native Testing Library)

```typescript
// __tests__/SubscriptionCard.test.tsx
describe('SubscriptionCard', () => {
  it('renders subscription name and price', () => {
    const sub = mockSubscription({ name: 'Netflix', price: 1490, currency: 'RUB' });
    const { getByText } = render(<SubscriptionCard subscription={sub} />);
    
    expect(getByText('Netflix')).toBeTruthy();
    expect(getByText('₽1,490')).toBeTruthy();
  });

  it('shows paused badge for paused subscriptions', () => {
    const sub = mockSubscription({ status: 'paused' });
    const { getByTestId } = render(<SubscriptionCard subscription={sub} />);
    expect(getByTestId('paused-badge')).toBeTruthy();
  });

  it('calls onPress when tapped', () => {
    const onPress = jest.fn();
    const { getByTestId } = render(
      <SubscriptionCard subscription={mockSubscription()} onPress={onPress} />
    );
    fireEvent.press(getByTestId('subscription-card'));
    expect(onPress).toHaveBeenCalledTimes(1);
  });
});
```

#### 3.2 E2E-тесты (Detox)

```javascript
// e2e/addSubscription.test.js
describe('Add Subscription Flow', () => {
  beforeAll(async () => {
    await device.launchApp();
    await loginAs('test@example.com');
  });

  it('should add a subscription manually', async () => {
    await element(by.id('tab-subscriptions')).tap();
    await element(by.id('btn-add-subscription')).tap();
    
    await element(by.id('input-name')).typeText('Netflix');
    await element(by.id('input-price')).typeText('1490');
    await element(by.id('select-billing-cycle')).tap();
    await element(by.text('Ежемесячно')).tap();
    await element(by.id('btn-save')).tap();
    
    await expect(element(by.text('Netflix'))).toBeVisible();
  });
});
```

#### 3.3 Чек-лист тестирования

| Сценарий | Тип | Статус |
|----------|-----|--------|
| SubscriptionCard renders correctly | Unit | ⬜ |
| Dashboard shows total amount | Unit | ⬜ |
| Add subscription form validation | Unit | ⬜ |
| Full add subscription E2E flow | Detox | ⬜ |
| Login with biometrics | Detox | ⬜ |
| Pull-to-refresh updates list | Detox | ⬜ |
| Offline mode shows cached data | Detox | ⬜ |
| Push notification deep link works | Detox | ⬜ |
| Performance: список 100 подписок < 100ms render | Perf | ⬜ |
| Тёмная тема: все компоненты корректны | Visual | ⬜ |
| iOS: работает на iPhone 14 Pro | Device | ⬜ |
| Android: работает на Samsung Galaxy S22 | Device | ⬜ |

---

## Модуль 14: Веб-приложение

### Фаза 1: Планирование задач

**Длительность:** 3 дня  
**Участники:** Frontend Developer ×2, UI/UX Designer, Tech Lead

#### 1.1 Задачи планирования

- [ ] Определить страницы App Router: layout, page, loading, error
- [ ] Согласовать использование Server Components vs Client Components
- [ ] Спроектировать систему роутинга и защищённых роутов
- [ ] Выбрать библиотеку иконок: Lucide React
- [ ] Определить breakpoints для адаптивной вёрстки
- [ ] Настроить общий HTTP-клиент с interceptors
- [ ] Создать тикеты WEB-001..020

#### 1.2 Структура App Router

```
app/
├── (auth)/
│   ├── login/
│   │   └── page.tsx
│   ├── register/
│   │   └── page.tsx
│   └── forgot-password/
│       └── page.tsx
├── (dashboard)/
│   ├── layout.tsx          # Sidebar + Header (Server Component)
│   ├── page.tsx            # Dashboard
│   ├── subscriptions/
│   │   ├── page.tsx        # Список
│   │   ├── [id]/page.tsx   # Детали
│   │   └── new/page.tsx    # Форма
│   ├── analytics/
│   │   └── page.tsx
│   ├── notifications/
│   │   └── page.tsx
│   ├── integrations/
│   │   └── page.tsx
│   └── settings/
│       └── page.tsx
├── (marketing)/
│   ├── page.tsx            # Лендинг
│   ├── pricing/page.tsx
│   └── about/page.tsx
└── api/
    └── auth/[...nextauth]/route.ts
```

#### 1.3 Ключевые решения архитектуры

```
Server Components:    Layout, Sidebar, статичные части страниц
Client Components:    Интерактивные формы, графики, real-time данные
Data Fetching:        TanStack Query (client) + fetch с кэшем (server)
Auth:                 NextAuth.js (JWT strategy, кастомные провайдеры)
Shared Code с RN:     ~/shared/ — типы TypeScript, утилиты, константы
```

---

### Фаза 2: Разработка

#### 2.1 WEB-006: Dashboard страница

```tsx
// app/(dashboard)/page.tsx — Server Component

export default async function DashboardPage() {
  // Серверный prefetch данных (SSR)
  const session = await getServerSession();
  const overview = await fetchOverview(session.token);
  const upcoming = await fetchUpcoming(session.token);

  return (
    <div className="space-y-6">
      <OverviewCards initialData={overview} />       {/* Client */}
      <div className="grid grid-cols-2 gap-6">
        <SpendingChart />                             {/* Client */}
        <CategoryPieChart />                          {/* Client */}
      </div>
      <UpcomingBillings initialData={upcoming} />    {/* Client */}
    </div>
  );
}
```

#### 2.2 WEB-010: Страница аналитики (Recharts)

```tsx
// components/analytics/SpendingChart.tsx
'use client';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';

export const SpendingChart = () => {
  const { data } = useQuery(['analytics', 'monthly'], fetchMonthlySpending);

  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={data?.months}>
        <XAxis dataKey="month" />
        <YAxis tickFormatter={(v) => `₽${v.toLocaleString()}`} />
        <Tooltip formatter={(v) => [`₽${Number(v).toLocaleString()}`, 'Затраты']} />
        <Bar dataKey="total" fill="#6366f1" radius={[4, 4, 0, 0]} />
      </BarChart>
    </ResponsiveContainer>
  );
};
```

---

### Фаза 3: Тестирование

#### 3.1 Unit/Component тесты (Jest + RTL)

```typescript
// __tests__/SubscriptionsTable.test.tsx
describe('SubscriptionsTable', () => {
  it('renders all subscriptions', () => {
    const subs = [mockSub({ name: 'Netflix' }), mockSub({ name: 'Spotify' })];
    render(<SubscriptionsTable subscriptions={subs} />);
    expect(screen.getByText('Netflix')).toBeInTheDocument();
    expect(screen.getByText('Spotify')).toBeInTheDocument();
  });

  it('filters by active status', async () => {
    render(<SubscriptionsTable subscriptions={mockSubs} />);
    await userEvent.click(screen.getByRole('button', { name: 'Активные' }));
    expect(screen.queryByText('Paused Sub')).not.toBeInTheDocument();
  });

  it('opens edit form on row click', async () => {
    render(<SubscriptionsTable subscriptions={[mockSub()]} />);
    await userEvent.click(screen.getByTestId('subscription-row'));
    expect(screen.getByRole('dialog')).toBeInTheDocument();
  });
});
```

#### 3.2 E2E тесты (Playwright)

```typescript
// e2e/subscriptions.spec.ts
test('add subscription flow', async ({ page }) => {
  await loginAs(page, 'test@example.com');
  await page.goto('/subscriptions/new');
  
  await page.fill('[name="name"]', 'Netflix');
  await page.fill('[name="price"]', '1490');
  await page.selectOption('[name="billing_cycle"]', 'monthly');
  await page.fill('[name="next_billing_date"]', '2026-03-15');
  await page.click('button[type="submit"]');
  
  await expect(page.locator('text=Netflix')).toBeVisible();
  await expect(page.locator('text=₽1,490')).toBeVisible();
});

test('analytics page shows charts', async ({ page }) => {
  await loginAs(page, 'test@example.com');
  await page.goto('/analytics');
  
  await expect(page.locator('[data-testid="spending-chart"]')).toBeVisible();
  await expect(page.locator('[data-testid="category-pie"]')).toBeVisible();
});
```

#### 3.3 Performance-тесты (Lighthouse CI)

```yaml
# .lighthouserc.yml
ci:
  collect:
    url:
      - https://staging.multisub.app/
      - https://staging.multisub.app/subscriptions
      - https://staging.multisub.app/analytics
  assert:
    assertions:
      first-contentful-paint:     ['error', {maxNumericValue: 1800}]
      largest-contentful-paint:   ['error', {maxNumericValue: 2500}]
      total-blocking-time:        ['warn',  {maxNumericValue: 300}]
      cumulative-layout-shift:    ['error', {maxNumericValue: 0.1}]
      interactive:               ['error', {maxNumericValue: 3500}]
```

#### 3.4 Чек-лист тестирования

| Сценарий | Тип | Статус |
|----------|-----|--------|
| Таблица отображает все подписки | Unit | ⬜ |
| Фильтр по статусу | Unit | ⬜ |
| Форма добавления: валидация | Unit | ⬜ |
| Dashboard SSR: данные на первом ренере | Integration | ⬜ |
| E2E: добавление подписки | Playwright | ⬜ |
| E2E: логин → дашборд | Playwright | ⬜ |
| E2E: экспорт CSV | Playwright | ⬜ |
| LCP < 2.5s (Lighthouse) | Performance | ⬜ |
| CLS < 0.1 (Lighthouse) | Performance | ⬜ |
| Accessibility: keyboard navigation | Accessibility | ⬜ |
| Accessibility: screen reader (NVDA) | Accessibility | ⬜ |
| Mobile: 375px viewport корректен | Responsive | ⬜ |
| Tablet: 768px viewport корректен | Responsive | ⬜ |

---

## 📊 Сводная матрица тестового покрытия

| Модуль | Unit | Integration | E2E | Load | Security |
|--------|:----:|:-----------:|:---:|:----:|:--------:|
| Auth | ✅ | ✅ | ✅ | ✅ | ✅ |
| User Profile | ✅ | ✅ | - | - | - |
| Subscriptions | ✅ | ✅ | ✅ | - | - |
| Bank Integration | ✅ | ✅ | - | - | ✅ |
| Email Parsing | ✅ | ✅ | - | - | ✅ |
| Analytics | ✅ | ✅ | - | ✅ | - |
| Usage Tracking | ✅ | ✅ | - | - | - |
| Prediction | ✅ | ✅ | - | - | - |
| Notifications | ✅ | ✅ | - | ✅ | - |
| Catalog | ✅ | ✅ | - | - | - |
| Cancellation | ✅ | ✅ | - | - | - |
| Admin | ✅ | ✅ | - | - | ✅ |
| Mobile App | ✅ | - | ✅ | ✅ | - |
| Web App | ✅ | ✅ | ✅ | ✅ | - |

---

## 🚀 CI/CD Pipeline (GitHub Actions — полный рабочий флоу)

### Структура репозитория

Репозиторий на GitHub организован по **monorepo** схеме:

```
multiSub/                       ← корень репозитория
├── .github/
│   ├── workflows/
│   │   ├── backend.yml        ← Go: lint + test + build + deploy (Railway)
│   │   ├── web.yml            ← Next.js: lint + test + Vercel deploy
│   │   ├── mobile.yml         ← RN: lint + test + EAS build
│   │   └── pr-check.yml       ← проверки для любого PR
│   └── PULL_REQUEST_TEMPLATE.md
├── skills/                     ← AI Agent скиллы (коммитятся в репо)
├── apps/
│   ├── backend/               ← Go API-сервер
│   ├── web/                   ← Next.js веб-приложение
│   └── mobile/                ← React Native (Expo)
├── packages/
│   ├── shared-types/          ← общие TypeScript типы
│   └── api-client/            ← сгенерированный API-клиент
├── docker-compose.yml         ← локальная разработка
├── TECHNICAL_SPECIFICATION.md
└── DEVELOPMENT_PLAN.md
```

---

### Стратегия веток (Git Flow для MVP)

```
main        ← • защищённая ветка
            • только через PR из develop
            • требует 1 review
            • деплой на production

develop     ← • ветка интеграции
            • деплой на staging
            • PR из feature-веток

feature/AUTH-001-registration
feature/SUB-001-create-subscription
fix/AUTH-003-token-refresh-bug
chore/update-dependencies
ci/add-telegram-notifications
```

**Формат commit-сообщений (Conventional Commits):**

```
feat(auth): add JWT refresh token rotation
fix(subscriptions): correct annual cost for quarterly billing
docs(api): update swagger spec for /subscriptions endpoint
test(auth): add brute-force protection tests
chore(deps): bump go to 1.22.3
ci: add telegram deployment notifications
```

---

### Workflow 1: `backend.yml` — Go API

```yaml
name: Backend CI/CD

on:
  push:
    branches: [main, develop]
    paths: ["apps/backend/**"]
  pull_request:
    branches: [main, develop]
    paths: ["apps/backend/**"]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}/backend

jobs:
  test:
    name: Lint & Test
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_DB: multisub_test
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
        ports: ["5432:5432"]
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache-dependency-path: apps/backend/go.sum
      - name: Install dependencies
        working-directory: apps/backend
        run: go mod download
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          working-directory: apps/backend
      - name: Run migrations
        working-directory: apps/backend
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/multisub_test?sslmode=disable
        run: go run ./cmd/migrate up
      - name: Run tests
        working-directory: apps/backend
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/multisub_test?sslmode=disable
          JWT_SECRET: ci-test-secret
        run: |
          go test ./... -race -count=1 \
            -coverprofile=coverage.out -covermode=atomic -timeout=120s
      - name: Check coverage >= 80%
        working-directory: apps/backend
        run: |
          COV=$(go tool cover -func coverage.out | grep total | awk '{print $3}' | tr -d '%')
          echo "Coverage: $COV%"
          awk "BEGIN { exit ($COV < 80) }"
      - uses: codecov/codecov-action@v4
        with:
          file: apps/backend/coverage.out
          flags: backend

  build:
    name: Build & Push Docker image
    runs-on: ubuntu-latest
    needs: test
    if: github.ref == 'refs/heads/main' || github.ref == 'refs/heads/develop'
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v4
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/metadata-action@v5
        id: meta
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=sha,prefix={{branch}}-
            type=ref,event=branch
      - uses: docker/build-push-action@v5
        with:
          context: apps/backend
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy-staging:
    name: Deploy → Staging (Railway)
    runs-on: ubuntu-latest
    needs: build
    if: github.ref == 'refs/heads/develop'
    environment: staging
    steps:
      - uses: actions/checkout@v4
      - run: npm install -g @railway/cli
      - name: Deploy
        env:
          RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN_STAGING }}
        run: railway up --service backend --environment staging
      - name: Smoke test
        run: sleep 15 && curl -f https://api-staging.multisub.app/health
      - name: Notify Telegram
        if: always()
        uses: appleboy/telegram-action@master
        with:
          to: ${{ secrets.TELEGRAM_CHAT_ID }}
          token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
          message: |
            ${{ job.status == 'success' && '✅' || '❌' }} Backend Staging
            Branch: ${{ github.ref_name }} | By: ${{ github.actor }}

  deploy-production:
    name: Deploy → Production (Railway)
    runs-on: ubuntu-latest
    needs: build
    if: github.ref == 'refs/heads/main'
    environment: production
    steps:
      - uses: actions/checkout@v4
      - run: npm install -g @railway/cli
      - name: Deploy
        env:
          RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN_PROD }}
        run: railway up --service backend --environment production
      - name: Smoke test
        run: sleep 15 && curl -f https://api.multisub.app/health
```

---

### Workflow 2: `web.yml` — Next.js + Vercel

```yaml
name: Web CI/CD

on:
  push:
    branches: [main, develop]
    paths: ["apps/web/**", "packages/**"]
  pull_request:
    paths: ["apps/web/**", "packages/**"]

jobs:
  test:
    name: Lint & Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
      - run: npm ci
      - name: TypeScript check
        working-directory: apps/web
        run: npx tsc --noEmit
      - name: ESLint
        working-directory: apps/web
        run: npx eslint . --max-warnings 0
      - name: Jest tests
        working-directory: apps/web
        run: npx jest --coverage --ci
      - uses: codecov/codecov-action@v4
        with:
          flags: web

  deploy:
    name: Deploy → Vercel
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@v4
      - uses: amondnet/vercel-action@v25
        id: deploy
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
          vercel-project-id: ${{ secrets.VERCEL_PROJECT_ID }}
          # main → production preview, остальные → preview URL
          vercel-args: ${{ github.ref == 'refs/heads/main' && '--prod' || '' }}
          working-directory: apps/web
      # Комментарий в PR со ссылкой на preview
      - if: github.event_name == 'pull_request'
        uses: actions/github-script@v7
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: '🚀 Preview: ${{ steps.deploy.outputs.preview-url }}'
            })
```

---

### Workflow 3: `mobile.yml` — React Native (Expo EAS)

```yaml
name: Mobile CI/CD

on:
  push:
    branches: [main, develop]
    paths: ["apps/mobile/**", "packages/**"]
  pull_request:
    paths: ["apps/mobile/**"]

jobs:
  test:
    name: Lint & Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
      - run: npm ci
      - name: TypeScript check
        working-directory: apps/mobile
        run: npx tsc --noEmit
      - name: Jest tests
        working-directory: apps/mobile
        run: npx jest --coverage --ci

  eas-preview:
    name: EAS Build (Preview)
    runs-on: ubuntu-latest
    needs: test
    if: github.ref == 'refs/heads/develop'
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '20' }
      - run: npm install -g eas-cli
      - name: Build preview
        working-directory: apps/mobile
        env:
          EXPO_TOKEN: ${{ secrets.EXPO_TOKEN }}
        run: eas build --platform all --profile preview --non-interactive

  eas-production:
    name: EAS Build + Submit (Production)
    runs-on: ubuntu-latest
    needs: test
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '20' }
      - run: npm install -g eas-cli
      - name: Build & submit
        working-directory: apps/mobile
        env:
          EXPO_TOKEN: ${{ secrets.EXPO_TOKEN }}
        run: eas build --platform all --profile production --non-interactive
```

---

### Workflow 4: `pr-check.yml` — PR проверки

```yaml
name: PR Checks

on:
  pull_request:
    branches: [main, develop]

jobs:
  pr-title:
    name: Conventional Commits title check
    runs-on: ubuntu-latest
    steps:
      - uses: amannn/action-semantic-pull-request@v5
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          types: [feat, fix, docs, style, refactor, test, chore, ci, perf]
```

---

### GitHub Secrets — что настроить в репозитории

> `Settings → Secrets and variables → Actions → New repository secret`

| Secret | Где взять |
|--------|----------|
| `RAILWAY_TOKEN_STAGING` | Railway Dashboard → Account → Tokens |
| `RAILWAY_TOKEN_PROD` | Railway Dashboard → Account → Tokens |
| `VERCEL_TOKEN` | vercel.com → Settings → Tokens |
| `VERCEL_ORG_ID` | vercel.com → Settings → General |
| `VERCEL_PROJECT_ID` | vercel.com → Project → Settings |
| `EXPO_TOKEN` | expo.dev → Account → Access Tokens |
| `TELEGRAM_BOT_TOKEN` | @BotFather в Telegram |
| `TELEGRAM_CHAT_ID` | ID чата для уведомлений |

### PR Template (`.github/PULL_REQUEST_TEMPLATE.md`)

```markdown
## Описание
<!-- Что сделано? Зачем? -->

## Тип изменения
- [ ] feat (новая функция)
- [ ] fix (исправление бага)
- [ ] refactor / test / docs / chore / ci

## Связанные задачи
<!-- Closes #123 -->

## Чеклист
- [ ] Тесты написаны и проходят
- [ ] TypeScript / lint ошибок нет
- [ ] Migration добавлена (если нужно)
- [ ] Swagger/OpenAPI обновлён
- [ ] Нет breaking changes (или они задокументированы)
```

---

**Документ подготовлен:** MultiSub Team  
**Связан с:** [TECHNICAL_SPECIFICATION.md](./TECHNICAL_SPECIFICATION.md)  
**Версия:** 1.1 (MVP)
