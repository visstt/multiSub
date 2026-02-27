# Техническое задание: MultiSub - Система управления подписками

**Версия документа:** 1.1  
**Дата создания:** 27 февраля 2026  
**Статус:** MVP-редакция (бесплатный стек)

---

## 📋 Оглавление

1. [Общее описание проекта](#1-общее-описание-проекта)
2. [Стек технологий](#2-стек-технологий)
3. [Архитектура системы](#3-архитектура-системы)
4. [Модули системы](#4-модули-системы)
5. [План разработки](#5-план-разработки)
6. [Требования к безопасности](#6-требования-к-безопасности)
7. [Нефункциональные требования](#7-нефункциональные-требования)

---

## 1. Общее описание проекта

### 1.1 Цель проекта

Разработка кроссплатформенного решения **MultiSub** для централизованного управления всеми подписками пользователя с функциями автоматического сбора данных, аналитики, прогнозирования и умных уведомлений.

### 1.2 Целевая аудитория

- Пользователи с множественными подписками (от 3 и более)
- Люди, стремящиеся оптимизировать расходы
- Возраст: 18-55 лет
- География: РФ и СНГ (первичный рынок)

### 1.3 Платформы

| Платформа | Минимальная версия |
|-----------|-------------------|
| iOS | 14.0+ |
| Android | 8.0 (API 26)+ |
| Web | Chrome 90+, Firefox 88+, Safari 14+, Edge 90+ |

### 1.4 Основные функциональные возможности

1. **Единая панель управления** - отображение всех активных подписок
2. **Автоматический сбор данных** - интеграция с банками и почтой
3. **Аналитика затрат** - визуализация расходов по категориям
4. **Анализ использования** - отслеживание активности в подписках
5. **Прогнозирование** - расчёт затрат на год вперёд
6. **Умные уведомления** - предупреждения о списаниях и изменениях
7. **Поиск альтернатив** - рекомендации более выгодных сервисов
8. **Быстрая отмена** - помощь в отмене/приостановке подписок

---

## 2. Стек технологий

> ⚠️ **MVP-ограничение:** на этапе MVP используются исключительно бесплатные open-source технологии и бесплатные планы облачных платформ. Платные сервисы (Redis Cloud, MongoDB Atlas, Elasticsearch, Kubernetes) вводятся в Фазе 2–3 при масштабировании.

### 2.1 Backend (MVP — бесплатный стек)

| Компонент | Технология (MVP) | Обоснование |
|-----------|-----------------|-------------|
| **Язык** | Go 1.22+ | Высокая производительность, бесплатный open-source |
| **Веб-фреймворк** | Fiber v3 | Быстрый, минималистичный, бесплатный |
| **ORM / Query** | sqlc + goose | Генерация типизированного кода из SQL; goose — миграции (оба бесплатны) |
| **API** | REST (JSON) | Только REST на MVP; GraphQL в Фазе 2 |
| **Очереди (MVP)** | PostgreSQL LISTEN/NOTIFY + таблица jobs | Без дополнительных сервисов; бесплатно |
| **Кэширование (MVP)** | In-memory (go-cache / ristretto) | Бесплатно; Redis добавляется в Фазе 2 |
| **Поиск (MVP)** | PostgreSQL Full-Text Search (tsvector) | Встроен в PG; Elasticsearch — Фаза 2 |
| **Email** | resend.com (3 000 писем/мес бесплатно) | Бесплатный план покрывает MVP |
| **Хранилище файлов** | Cloudflare R2 (10 GB / мес бесплатно) | Аватары, вложения без платы за исходящий трафик |

### 2.2 Базы данных (MVP)

| База | Назначение | Хостинг (бесплатно) |
|------|------------|---------------------|
| **PostgreSQL 16+** | Единственная БД MVP: основные данные, очередь задач, FTS, аналитика | Railway (бесплатный план 1 GB) / Supabase (500 MB free tier) |

> **Убраны для MVP:** MongoDB, ClickHouse, Redis, Elasticsearch — слишком дорого и избыточно для первой версии. Добавляются только по необходимости при росте нагрузки.

### 2.3 Mobile (iOS & Android)

| Компонент | Технология | Обоснование |
|-----------|------------|-------------|
| **Фреймворк** | React Native 0.73+ | Единая кодовая база для iOS и Android, общий стек с вебом |
| **Язык** | TypeScript 5.3+ | Типизация, единый язык с веб-приложением |
| **State Management** | Zustand / TanStack Query | Лёгкое управление, общий подход с вебом |
| **Navigation** | React Navigation 6+ | Стандарт навигации в RN |
| **Networking** | Axios / Fetch API | HTTP-клиент с интерцепторами |
| **Local Storage** | MMKV / WatermelonDB | Высокопроизводительное локальное хранилище |
| **UI Components** | React Native Paper / NativeWind | Material Design + Tailwind-стиль стилизация |
| **Push** | Firebase Cloud Messaging | Кроссплатформенные пуши |

### 2.4 Web Application

| Компонент | Технология | Обоснование |
|-----------|------------|-------------|
| **Фреймворк** | Next.js 14+ | SSR, отличный DX, оптимизации |
| **Язык** | TypeScript 5.3+ | Типизация, надёжность кода |
| **UI Library** | React 18+ | Компонентный подход |
| **State** | Zustand / TanStack Query | Лёгкое управление состоянием + data fetching |
| **Styling** | Tailwind CSS 3.4+ | Utility-first подход |
| **Charts** | Recharts / Chart.js | Визуализация аналитики |
| **Forms** | React Hook Form + Zod | Формы с валидацией |

### 2.5 DevOps & Infrastructure (MVP — только бесплатное)

| Компонент | Технология MVP | Платформа / Цена |
|-----------|---------------|------------------|
| **Контейнеризация** | Docker + Docker Compose | Бесплатно, open-source |
| **Оркестрация** | — (нет K8s на MVP) | Один инстанс Railway/Render достаточен |
| **CI/CD** | **GitHub Actions** | Бесплатно: 2 000 мин/мес для приватных репо, ∞ для публичных |
| **Хостинг Web** | **Vercel** (free Hobby plan) | Автодеплой из GitHub, CDN, preview-ссылки для каждого PR |
| **Хостинг Backend** | **Railway** (free $5 credit/мес) | Docker-деплой, автодеплой из GitHub |
| **Хостинг DB** | **Railway PostgreSQL** или **Supabase** (free tier) | Встроенный в Railway / 500 MB Supabase |
| **Мониторинг** | Sentry (5 K ошибок/мес бесплатно) | Error tracking для backend и frontend |
| **Логирование** | Railway built-in logs + Axiom (бесплатный план) | Без ELK |
| **Uptime** | Better Uptime / UptimeRobot (бесплатный план) | Проверка доступности |
| **Мобильная сборка** | Expo EAS Build (30 сборок/мес бесплатно) | iOS + Android OTA-обновления |

#### Схема деплоя MVP

```
GitHub Repository
├── main         → автодеплой на Production
├── develop      → автодеплой на Staging
└── feature/*    → Preview Deploy (Vercel) + CI проверки

                          GitHub Actions
                         ┌──────────────────────────┐
  git push ─────────────▶│ 1. Lint + Tests          │
                         │ 2. Build Docker image    │
                         │ 3. Push to GHCR          │
                         │ 4. Deploy to Railway     │──▶ backend.multisub.app
                         │ 5. Deploy to Vercel      │──▶ multisub.app
                         │ 6. Notify Telegram       │
                         └──────────────────────────┘
```

#### GitHub репозиторий: структура веток

```
main          ← стабильный production, защищённая ветка (требует PR + 1 review)
develop       ← основная ветка разработки, автодеплой на staging
feature/XXX   ← задачи из Jira/Linear, merge в develop через PR
fix/XXX       ← хотфиксы, merge в main + develop
release/X.Y.Z ← release candidate перед merge в main
```

#### Настройки защиты main ветки

- Требуется 1 approved review
- CI должен пройти (lint, test, build)
- Нет прямых пушей (только через PR)
- Автоматическое удаление feature-веток после merge

### 2.6 Интеграции

| Сервис | Метод интеграции |
|--------|------------------|
| **Сбербанк Online** | Open Banking API (СБП) |
| **Тинькофф** | Tinkoff Invest API / Open API |
| **ЮMoney** | ЮMoney API |
| **Mail.ru** | IMAP + OAuth 2.0 |
| **Yandex Mail** | IMAP + OAuth 2.0 |
| **Gmail** | Gmail API + OAuth 2.0 |

---

## 3. Архитектура системы

### 3.1 Высокоуровневая архитектура

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              КЛИЕНТСКИЙ СЛОЙ                                │
├─────────────────┬─────────────────────────────────┬─────────────────────────┤
│   iOS App       │         Android App             │       Web App           │
│  (React Native) │        (React Native)            │       (Next.js)         │
└────────┬────────┴────────────────┬────────────────┴────────────┬────────────┘
         │                         │                              │
         └─────────────────────────┼──────────────────────────────┘
                                   │
                            ┌──────▼──────┐
                            │   API       │
                            │   Gateway   │
                            │   (Kong)    │
                            └──────┬──────┘
                                   │
┌──────────────────────────────────┼──────────────────────────────────────────┐
│                          BACKEND SERVICES (Микросервисы)                    │
├─────────────┬─────────────┬──────┴──────┬─────────────┬─────────────────────┤
│   Auth      │  Users      │ Subscriptions│ Analytics  │   Notifications     │
│   Service   │  Service    │   Service   │   Service   │     Service         │
├─────────────┼─────────────┼─────────────┼─────────────┼─────────────────────┤
│ Integration │ Parsing     │ Prediction  │ Alternative │   Template          │
│   Service   │ Service     │   Service   │   Service   │     Service         │
└──────┬──────┴──────┬──────┴──────┬──────┴──────┬──────┴──────────┬──────────┘
       │             │             │             │                 │
       └─────────────┴─────────────┼─────────────┴─────────────────┘
                                   │
┌──────────────────────────────────┼──────────────────────────────────────────┐
│                              DATA LAYER                                     │
├─────────────────┬────────────────┼────────────────┬─────────────────────────┤
│   PostgreSQL    │    MongoDB     │     Redis      │      ClickHouse         │
│   (Primary DB)  │    (Logs)      │    (Cache)     │      (Analytics)        │
└─────────────────┴────────────────┴────────────────┴─────────────────────────┘
```

### 3.2 Схема базы данных (PostgreSQL)

```sql
-- Основные таблицы

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    timezone VARCHAR(50) DEFAULT 'Europe/Moscow',
    language VARCHAR(10) DEFAULT 'ru',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    service_id UUID REFERENCES services(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    billing_cycle VARCHAR(20) NOT NULL, -- monthly, yearly, weekly
    next_billing_date DATE,
    start_date DATE,
    status VARCHAR(20) DEFAULT 'active', -- active, paused, cancelled
    auto_renew BOOLEAN DEFAULT true,
    source VARCHAR(50), -- manual, email_parse, bank_sync
    category_id UUID REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    logo_url VARCHAR(500),
    website_url VARCHAR(500),
    cancel_url VARCHAR(500),
    support_email VARCHAR(255),
    category_id UUID REFERENCES categories(id),
    avg_price DECIMAL(10, 2),
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    icon VARCHAR(50),
    parent_id UUID REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payment_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID REFERENCES subscriptions(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    payment_date TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL, -- success, failed, pending
    source VARCHAR(50), -- bank_sync, email_parse, manual
    raw_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE usage_tracking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID REFERENCES subscriptions(id) ON DELETE CASCADE,
    usage_date DATE NOT NULL,
    usage_type VARCHAR(50), -- login, content_view, feature_use
    usage_count INTEGER DEFAULT 1,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    subscription_id UUID REFERENCES subscriptions(id),
    type VARCHAR(50) NOT NULL, -- billing_reminder, price_change, inactivity
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    scheduled_at TIMESTAMP,
    sent_at TIMESTAMP,
    read_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE email_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- gmail, yandex, mailru
    email VARCHAR(255) NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP,
    last_sync_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bank_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    bank_name VARCHAR(100) NOT NULL,
    account_id VARCHAR(255),
    access_token TEXT ENCRYPTED,
    refresh_token TEXT ENCRYPTED,
    token_expires_at TIMESTAMP,
    last_sync_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notification_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    billing_reminder_days INTEGER[] DEFAULT '{1, 3, 7}',
    price_change_alert BOOLEAN DEFAULT true,
    inactivity_days INTEGER DEFAULT 30,
    email_notifications BOOLEAN DEFAULT true,
    push_notifications BOOLEAN DEFAULT true,
    quiet_hours_start TIME,
    quiet_hours_end TIME,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 4. Модули системы

---

### 📦 МОДУЛЬ 1: Аутентификация и авторизация (Auth Module)

**Описание:** Управление пользователями, регистрация, вход, восстановление пароля, OAuth.

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| AUTH-001 | Регистрация пользователя через email/пароль | High | 8 | - |
| AUTH-002 | Валидация email (отправка кода подтверждения) | High | 6 | AUTH-001 |
| AUTH-003 | Авторизация (логин) с JWT токенами | High | 8 | AUTH-001 |
| AUTH-004 | Refresh token механизм | High | 6 | AUTH-003 |
| AUTH-005 | Восстановление пароля через email | High | 6 | AUTH-002 |
| AUTH-006 | OAuth 2.0 через Google | Medium | 12 | AUTH-003 |
| AUTH-007 | OAuth 2.0 через Яндекс | Medium | 10 | AUTH-003 |
| AUTH-008 | OAuth 2.0 через VK ID | Medium | 10 | AUTH-003 |
| AUTH-009 | Двухфакторная аутентификация (TOTP) | Medium | 16 | AUTH-003 |
| AUTH-010 | Управление сессиями (просмотр, завершение) | Low | 8 | AUTH-003 |
| AUTH-011 | Rate limiting для защиты от брутфорса | High | 6 | - |
| AUTH-012 | Логирование попыток входа | Medium | 4 | AUTH-003 |

**API Endpoints:**

```yaml
POST   /api/v1/auth/register          # Регистрация
POST   /api/v1/auth/login              # Вход
POST   /api/v1/auth/logout             # Выход
POST   /api/v1/auth/refresh            # Обновление токена
POST   /api/v1/auth/forgot-password    # Запрос сброса пароля
POST   /api/v1/auth/reset-password     # Сброс пароля
POST   /api/v1/auth/verify-email       # Подтверждение email
GET    /api/v1/auth/oauth/google       # OAuth Google
GET    /api/v1/auth/oauth/yandex       # OAuth Яндекс
GET    /api/v1/auth/oauth/vk           # OAuth VK
POST   /api/v1/auth/2fa/enable         # Включение 2FA
POST   /api/v1/auth/2fa/verify         # Проверка 2FA
GET    /api/v1/auth/sessions           # Список сессий
DELETE /api/v1/auth/sessions/:id       # Завершение сессии
```

---

### 📦 МОДУЛЬ 2: Профиль пользователя (User Profile Module)

**Описание:** Управление профилем, настройками, персональными данными.

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| USER-001 | Просмотр профиля пользователя | High | 4 | AUTH-003 |
| USER-002 | Редактирование профиля (имя, фото, контакты) | High | 8 | USER-001 |
| USER-003 | Загрузка и обработка аватара | Medium | 8 | USER-002 |
| USER-004 | Изменение пароля | High | 4 | AUTH-003 |
| USER-005 | Изменение email с подтверждением | Medium | 8 | AUTH-002 |
| USER-006 | Настройки часового пояса и языка | Medium | 4 | USER-001 |
| USER-007 | Экспорт данных пользователя (GDPR) | Medium | 12 | USER-001 |
| USER-008 | Удаление аккаунта | Medium | 8 | USER-001 |
| USER-009 | История изменений профиля | Low | 6 | USER-002 |

**API Endpoints:**

```yaml
GET    /api/v1/users/me                # Получение профиля
PATCH  /api/v1/users/me                # Обновление профиля
POST   /api/v1/users/me/avatar         # Загрузка аватара
DELETE /api/v1/users/me/avatar         # Удаление аватара
PUT    /api/v1/users/me/password       # Изменение пароля
PUT    /api/v1/users/me/email          # Изменение email
GET    /api/v1/users/me/export         # Экспорт данных
DELETE /api/v1/users/me                # Удаление аккаунта
```

---

### 📦 МОДУЛЬ 3: Управление подписками (Subscriptions Module)

**Описание:** CRUD операции над подписками, категоризация, статусы.

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| SUB-001 | Создание подписки вручную | High | 8 | AUTH-003 |
| SUB-002 | Просмотр списка всех подписок | High | 6 | SUB-001 |
| SUB-003 | Детальный просмотр подписки | High | 6 | SUB-001 |
| SUB-004 | Редактирование подписки | High | 8 | SUB-001 |
| SUB-005 | Удаление/архивация подписки | High | 4 | SUB-001 |
| SUB-006 | Фильтрация подписок (по статусу, категории) | Medium | 8 | SUB-002 |
| SUB-007 | Сортировка подписок | Medium | 4 | SUB-002 |
| SUB-008 | Поиск по подпискам | Medium | 6 | SUB-002 |
| SUB-009 | Управление категориями подписок | Medium | 8 | SUB-001 |
| SUB-010 | Привязка подписки к сервису из каталога | Medium | 6 | SUB-001, CAT-001 |
| SUB-011 | Массовые операции (пауза, удаление) | Low | 8 | SUB-001 |
| SUB-012 | Дублирование подписки | Low | 4 | SUB-001 |
| SUB-013 | История изменений подписки | Low | 8 | SUB-001 |
| SUB-014 | Теги для подписок | Low | 6 | SUB-001 |

**API Endpoints:**

```yaml
GET    /api/v1/subscriptions                    # Список подписок
POST   /api/v1/subscriptions                    # Создание подписки
GET    /api/v1/subscriptions/:id                # Детали подписки
PATCH  /api/v1/subscriptions/:id                # Обновление подписки
DELETE /api/v1/subscriptions/:id                # Удаление подписки
POST   /api/v1/subscriptions/:id/pause          # Приостановка
POST   /api/v1/subscriptions/:id/resume         # Возобновление
POST   /api/v1/subscriptions/:id/duplicate      # Дублирование
GET    /api/v1/subscriptions/:id/history        # История изменений
POST   /api/v1/subscriptions/bulk               # Массовые операции
GET    /api/v1/categories                       # Список категорий
```

**Статусы подписок:**

```
┌──────────┐     pause      ┌──────────┐
│  ACTIVE  │ ────────────▶ │  PAUSED  │
│          │ ◀──────────── │          │
└────┬─────┘     resume     └──────────┘
     │
     │ cancel
     ▼
┌──────────┐                ┌──────────┐
│ CANCELLED│                │ EXPIRED  │
└──────────┘                └──────────┘
```

---

### 📦 МОДУЛЬ 4: Интеграция с банками (Bank Integration Module)

**Описание:** Подключение к банковским API для автоматического получения транзакций.

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| BANK-001 | OAuth авторизация в Сбербанк Online | High | 24 | AUTH-003 |
| BANK-002 | OAuth авторизация в Тинькофф | High | 24 | AUTH-003 |
| BANK-003 | OAuth авторизация в ЮMoney | High | 16 | AUTH-003 |
| BANK-004 | Получение списка транзакций | High | 16 | BANK-001/002/003 |
| BANK-005 | Фильтрация транзакций по паттернам подписок | High | 20 | BANK-004 |
| BANK-006 | Автоматическое сопоставление транзакций с подписками | High | 24 | BANK-005, SUB-001 |
| BANK-007 | Ручное подтверждение сопоставления | Medium | 8 | BANK-006 |
| BANK-008 | Периодическая синхронизация (scheduler) | High | 12 | BANK-004 |
| BANK-009 | Обработка ошибок и переподключение | High | 16 | BANK-001/002/003 |
| BANK-010 | Шифрование и безопасное хранение токенов | Critical | 16 | BANK-001/002/003 |
| BANK-011 | Уведомление о необходимости переподключения | Medium | 6 | BANK-009 |
| BANK-012 | Отключение банковского подключения | Medium | 4 | BANK-001/002/003 |
| BANK-013 | Dashboard статуса подключений | Medium | 8 | BANK-001/002/003 |

**API Endpoints:**

```yaml
GET    /api/v1/integrations/banks               # Список поддерживаемых банков
POST   /api/v1/integrations/banks/:bank/connect # Начало подключения
GET    /api/v1/integrations/banks/:bank/callback# OAuth callback
GET    /api/v1/integrations/banks/connections   # Активные подключения
DELETE /api/v1/integrations/banks/connections/:id # Отключение
POST   /api/v1/integrations/banks/connections/:id/sync # Ручная синхронизация
GET    /api/v1/integrations/banks/transactions  # Найденные транзакции
POST   /api/v1/integrations/banks/transactions/:id/link # Привязка к подписке
```

---

### 📦 МОДУЛЬ 5: Парсинг почты (Email Parsing Module)

**Описание:** Подключение почтовых аккаунтов и парсинг писем с чеками/уведомлениями.

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| EMAIL-001 | OAuth авторизация Gmail | High | 16 | AUTH-003 |
| EMAIL-002 | OAuth авторизация Яндекс.Почта | High | 16 | AUTH-003 |
| EMAIL-003 | OAuth авторизация Mail.ru | High | 16 | AUTH-003 |
| EMAIL-004 | IMAP подключение с токенами | High | 12 | EMAIL-001/002/003 |
| EMAIL-005 | Получение списка писем (с пагинацией) | High | 8 | EMAIL-004 |
| EMAIL-006 | Фильтрация писем по отправителям-сервисам | High | 12 | EMAIL-005 |
| EMAIL-007 | Парсинг HTML писем | High | 24 | EMAIL-005 |
| EMAIL-008 | Распознавание шаблонов чеков | High | 32 | EMAIL-007 |
| EMAIL-009 | Извлечение данных: сумма, дата, сервис | High | 24 | EMAIL-008 |
| EMAIL-010 | Машинное обучение для улучшения распознавания | Medium | 40 | EMAIL-008 |
| EMAIL-011 | Создание подписок из распознанных данных | High | 16 | EMAIL-009, SUB-001 |
| EMAIL-012 | Ручное подтверждение/корректировка | Medium | 8 | EMAIL-011 |
| EMAIL-013 | Периодическая проверка новых писем | High | 12 | EMAIL-005 |
| EMAIL-014 | Уведомление о найденных подписках | Medium | 6 | EMAIL-011 |
| EMAIL-015 | Управление правилами парсинга | Low | 16 | EMAIL-008 |
| EMAIL-016 | Отключение почтового подключения | Medium | 4 | EMAIL-001/002/003 |

**Поддерживаемые форматы чеков:**

- Стриминговые сервисы: Netflix, Spotify, YouTube Premium, Кинопоиск, Okko, ivi
- Софт: Adobe, Microsoft 365, JetBrains, Figma
- Облака: iCloud, Google One, Яндекс.Диск, Dropbox
- Доставка: Яндекс.Плюс, СберПрайм, Wildberries Premium
- И другие (расширяемый список)

**API Endpoints:**

```yaml
GET    /api/v1/integrations/email/providers     # Поддерживаемые провайдеры
POST   /api/v1/integrations/email/:provider/connect # Подключение
GET    /api/v1/integrations/email/:provider/callback # OAuth callback
GET    /api/v1/integrations/email/connections   # Активные подключения
DELETE /api/v1/integrations/email/connections/:id # Отключение
POST   /api/v1/integrations/email/connections/:id/sync # Синхронизация
GET    /api/v1/integrations/email/parsed        # Распознанные подписки
POST   /api/v1/integrations/email/parsed/:id/confirm # Подтверждение
POST   /api/v1/integrations/email/parsed/:id/reject  # Отклонение
```

---

### 📦 МОДУЛЬ 6: Аналитика (Analytics Module)

**Описание:** Визуализация затрат, статистика, группировка по категориям и периодам.

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| ANAL-001 | Общая статистика (сумма, кол-во подписок) | High | 8 | SUB-001 |
| ANAL-002 | График затрат по месяцам | High | 12 | ANAL-001 |
| ANAL-003 | Распределение по категориям (pie chart) | High | 10 | ANAL-001 |
| ANAL-004 | Топ-5 самых дорогих подписок | Medium | 6 | ANAL-001 |
| ANAL-005 | Сравнение месяц к месяцу | Medium | 10 | ANAL-002 |
| ANAL-006 | Сравнение год к году | Medium | 10 | ANAL-002 |
| ANAL-007 | Фильтры по периоду (месяц, квартал, год, custom) | High | 8 | ANAL-001 |
| ANAL-008 | Экспорт отчётов в PDF | Medium | 16 | ANAL-001..006 |
| ANAL-009 | Экспорт данных в CSV/Excel | Medium | 12 | ANAL-001 |
| ANAL-010 | Виджеты для dashboard | Medium | 16 | ANAL-001..006 |
| ANAL-011 | Кэширование агрегатов | High | 8 | ANAL-001 |
| ANAL-012 | Real-time обновление при изменениях | Low | 12 | ANAL-011 |

**API Endpoints:**

```yaml
GET    /api/v1/analytics/overview               # Общая статистика
GET    /api/v1/analytics/spending/monthly       # По месяцам
GET    /api/v1/analytics/spending/categories    # По категориям
GET    /api/v1/analytics/spending/services      # По сервисам
GET    /api/v1/analytics/top-subscriptions      # Топ подписок
GET    /api/v1/analytics/comparison             # Сравнение периодов
GET    /api/v1/analytics/export/pdf             # Экспорт PDF
GET    /api/v1/analytics/export/csv             # Экспорт CSV
```

**Виджеты Dashboard:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          MONTHLY OVERVIEW                               │
├──────────────┬──────────────┬──────────────┬───────────────────────────┤
│  Всего       │  Активных    │  Средняя     │     Изменение             │
│  ₽ 12,450    │  12 подписок │  ₽ 1,037     │     +5.2% ↑               │
├──────────────┴──────────────┴──────────────┴───────────────────────────┤
│                                                                         │
│    ▓▓▓▓▓▓                                                               │
│    ▓▓▓▓▓▓  ▓▓▓▓                                                         │
│    ▓▓▓▓▓▓  ▓▓▓▓  ▓▓▓                                                    │
│    ▓▓▓▓▓▓  ▓▓▓▓  ▓▓▓  ▓▓                                                │
│    ─────────────────────────                                            │
│    Янв  Фев  Мар  Апр  ...                                              │
│                                                                         │
├─────────────────────────────┬───────────────────────────────────────────┤
│   КАТЕГОРИИ                 │   ПРЕДСТОЯЩИЕ СПИСАНИЯ                   │
│   ● Стриминг    45%         │   • Netflix     ₽1,490  через 3 дня       │
│   ● Софт        25%         │   • Spotify     ₽169    через 5 дней      │
│   ● Облако      15%         │   • iCloud      ₽599    через 12 дней     │
│   ● Другое      15%         │                                           │
└─────────────────────────────┴───────────────────────────────────────────┘
```

---

### 📦 МОДУЛЬ 7: Отслеживание использования (Usage Tracking Module)

**Описание:** Учёт использования подписок пользователем.

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| USE-001 | Ручная отметка об использовании | High | 8 | SUB-001 |
| USE-002 | Quick-action кнопка "Использовал сегодня" | High | 4 | USE-001 |
| USE-003 | Календарь использования | Medium | 12 | USE-001 |
| USE-004 | Статистика использования по подписке | Medium | 10 | USE-001 |
| USE-005 | Расчёт "стоимости за использование" | Medium | 8 | USE-004 |
| USE-006 | Виджет неиспользуемых подписок | High | 8 | USE-004 |
| USE-007 | Напоминание об использовании | Medium | 6 | USE-006, NOTIF-001 |
| USE-008 | Интеграция с Screen Time (iOS) | Low | 20 | - |
| USE-009 | Интеграция с Digital Wellbeing (Android) | Low | 20 | - |
| USE-010 | Автоматическое отслеживание (опционально) | Low | 32 | USE-008/009 |

**API Endpoints:**

```yaml
POST   /api/v1/subscriptions/:id/usage          # Отметка использования
GET    /api/v1/subscriptions/:id/usage          # История использования
GET    /api/v1/subscriptions/:id/usage/stats    # Статистика
GET    /api/v1/usage/calendar                   # Календарь по всем подпискам
GET    /api/v1/usage/unused                     # Неиспользуемые подписки
DELETE /api/v1/subscriptions/:id/usage/:date    # Удаление отметки
```

---

### 📦 МОДУЛЬ 8: Прогнозирование (Prediction Module)

**Описание:** ML-модель для прогноза затрат на подписки.

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| PRED-001 | Базовый расчёт прогноза (сумма * 12) | High | 8 | SUB-001 |
| PRED-002 | Учёт разных биллинговых циклов | High | 12 | PRED-001 |
| PRED-003 | Учёт исторических изменений цен | Medium | 16 | PRED-001 |
| PRED-004 | Прогноз с учётом инфляции/трендов | Medium | 20 | PRED-003 |
| PRED-005 | Сценарии "что если" (добавить/убрать подписку) | Medium | 16 | PRED-001 |
| PRED-006 | Визуализация прогноза (график на год) | High | 12 | PRED-001 |
| PRED-007 | Помесячная детализация прогноза | Medium | 8 | PRED-006 |
| PRED-008 | Экспорт прогноза | Low | 6 | PRED-006 |
| PRED-009 | Уведомление о превышении бюджета | Medium | 8 | PRED-001, NOTIF-001 |
| PRED-010 | Установка бюджета на подписки | Medium | 6 | - |

**API Endpoints:**

```yaml
GET    /api/v1/predictions/yearly               # Годовой прогноз
GET    /api/v1/predictions/monthly              # Помесячный прогноз
POST   /api/v1/predictions/simulate             # Симуляция сценария
GET    /api/v1/predictions/budget               # Текущий бюджет
PUT    /api/v1/predictions/budget               # Установка бюджета
```

**Алгоритм прогнозирования:**

```
1. Собрать все активные подписки
2. Для каждой подписки:
   a. Определить базовую годовую стоимость
   b. Применить коэффициент изменения цены (историческая тенденция)
   c. Учесть вероятность отмены (на основе использования)
3. Агрегировать в годовой/месячный прогноз
4. Добавить доверительный интервал (±5-15%)
```

---

### 📦 МОДУЛЬ 9: Уведомления (Notifications Module)

**Описание:** Система умных уведомлений с различными триггерами.

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| NOTIF-001 | Инфраструктура push-уведомлений (FCM) | High | 16 | - |
| NOTIF-002 | Уведомления по email | High | 12 | - |
| NOTIF-003 | In-app уведомления | High | 10 | - |
| NOTIF-004 | Триггер: предстоящее списание (за N дней) | High | 12 | SUB-001 |
| NOTIF-005 | Триггер: изменение цены | High | 10 | SUB-001 |
| NOTIF-006 | Триггер: неактивность в сервисе | Medium | 12 | USE-001 |
| NOTIF-007 | Триггер: истечение пробного периода | Medium | 8 | SUB-001 |
| NOTIF-008 | Триггер: найдена более выгодная альтернатива | Low | 10 | ALT-001 |
| NOTIF-009 | Настройки уведомлений (типы, время) | High | 12 | NOTIF-001 |
| NOTIF-010 | Тихие часы (не беспокоить) | Medium | 6 | NOTIF-009 |
| NOTIF-011 | История уведомлений | Medium | 8 | NOTIF-001 |
| NOTIF-012 | Отметка прочитано/не прочитано | Medium | 4 | NOTIF-011 |
| NOTIF-013 | Массовое управление уведомлениями | Low | 6 | NOTIF-011 |
| NOTIF-014 | Scheduler для отложенных уведомлений | High | 16 | NOTIF-001 |
| NOTIF-015 | A/B тестирование текстов уведомлений | Low | 20 | NOTIF-001 |

**API Endpoints:**

```yaml
GET    /api/v1/notifications                    # Список уведомлений
GET    /api/v1/notifications/unread/count       # Счётчик непрочитанных
POST   /api/v1/notifications/:id/read           # Отметить прочитанным
POST   /api/v1/notifications/read-all           # Прочитать все
DELETE /api/v1/notifications/:id                # Удалить
GET    /api/v1/notifications/settings           # Настройки
PUT    /api/v1/notifications/settings           # Обновить настройки
POST   /api/v1/notifications/test               # Тестовое уведомление
```

**Типы уведомлений:**

| Тип | Описание | Приоритет | Канал по умолчанию |
|-----|----------|-----------|-------------------|
| `billing_reminder` | Напоминание о списании | High | Push + Email |
| `price_change` | Изменение цены | High | Push + Email |
| `trial_ending` | Окончание пробного периода | High | Push |
| `inactivity_alert` | Не используете подписку | Medium | Push |
| `alternative_found` | Найдена альтернатива | Low | In-app |
| `sync_required` | Требуется переподключение | High | Push |
| `weekly_summary` | Еженедельный дайджест | Low | Email |

---

### 📦 МОДУЛЬ 10: Каталог сервисов и альтернативы (Catalog & Alternatives Module)

**Описание:** База сервисов с возможностью поиска альтернатив.

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| CAT-001 | База данных популярных сервисов (~500+) | High | 40 | - |
| CAT-002 | CRUD для администрирования каталога | Medium | 16 | CAT-001 |
| CAT-003 | Поиск по каталогу (название, категория) | High | 12 | CAT-001 |
| CAT-004 | Детальная страница сервиса | High | 10 | CAT-001 |
| CAT-005 | Сравнение сервисов | Medium | 16 | CAT-001 |
| ALT-001 | Алгоритм подбора альтернатив | High | 24 | CAT-001 |
| ALT-002 | Фильтры альтернатив (цена, функции) | Medium | 12 | ALT-001 |
| ALT-003 | Рейтинг альтернатив (пользовательский) | Low | 16 | ALT-001 |
| ALT-004 | Интеграция с внешними API (цены, рейтинги) | Low | 24 | CAT-001 |
| ALT-005 | Персонализированные рекомендации | Medium | 20 | ALT-001, SUB-001 |
| ALT-006 | "Сэкономлено при переходе" калькулятор | Medium | 8 | ALT-001 |

**API Endpoints:**

```yaml
GET    /api/v1/catalog/services                 # Список сервисов
GET    /api/v1/catalog/services/:id             # Детали сервиса
GET    /api/v1/catalog/services/search          # Поиск
GET    /api/v1/catalog/categories               # Категории
GET    /api/v1/catalog/services/:id/alternatives# Альтернативы сервиса
POST   /api/v1/catalog/services/compare         # Сравнение сервисов
GET    /api/v1/recommendations                  # Персональные рекомендации
```

**Структура данных сервиса:**

```json
{
  "id": "uuid",
  "name": "Netflix",
  "description": "Стриминговый сервис фильмов и сериалов",
  "logo_url": "https://...",
  "website_url": "https://netflix.com",
  "category": {
    "id": "uuid",
    "name": "Видео стриминг",
    "slug": "video-streaming"
  },
  "pricing": [
    {"plan": "Базовый", "price": 799, "currency": "RUB", "billing_cycle": "monthly"},
    {"plan": "Стандарт", "price": 1190, "currency": "RUB", "billing_cycle": "monthly"},
    {"plan": "Премиум", "price": 1490, "currency": "RUB", "billing_cycle": "monthly"}
  ],
  "features": ["4K", "HDR", "Скачивание", "Несколько профилей"],
  "rating": 4.5,
  "reviews_count": 12500,
  "cancel_url": "https://netflix.com/cancelplan",
  "support_email": "support@netflix.com"
}
```

---

### 📦 МОДУЛЬ 11: Отмена подписок (Cancellation Module)

**Описание:** Помощь в отмене или приостановке подписок.

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| CANC-001 | Инструкции по отмене для каждого сервиса | High | 40 | CAT-001 |
| CANC-002 | Quick-action редирект на страницу отмены | High | 8 | CAT-001 |
| CANC-003 | Генератор письма в поддержку (шаблоны) | High | 16 | - |
| CANC-004 | Кастомизация шаблона письма | Medium | 8 | CANC-003 |
| CANC-005 | Копирование/отправка письма по email | Medium | 6 | CANC-003 |
| CANC-006 | Чек-лист "что сделать перед отменой" | Medium | 12 | - |
| CANC-007 | Отслеживание статуса отмены | Medium | 10 | SUB-001 |
| CANC-008 | Уведомление об успешной отмене | Medium | 6 | CANC-007, NOTIF-001 |
| CANC-009 | "Заморозка" вместо отмены (рекомендация) | Low | 8 | CANC-001 |
| CANC-010 | Сбор обратной связи (почему отменяете) | Low | 10 | CANC-007 |

**API Endpoints:**

```yaml
GET    /api/v1/subscriptions/:id/cancel-guide   # Инструкция по отмене
GET    /api/v1/subscriptions/:id/cancel-letter  # Шаблон письма
POST   /api/v1/subscriptions/:id/cancel-letter/customize # Кастомизация
POST   /api/v1/subscriptions/:id/cancel-letter/send # Отправка (через наш сервис)
GET    /api/v1/subscriptions/:id/pre-cancel-checklist # Чек-лист
POST   /api/v1/subscriptions/:id/mark-cancelled # Отметить как отменённую
```

**Шаблон письма для отмены:**

```
Тема: Запрос на отмену подписки - [Имя пользователя]

Здравствуйте!

Прошу отменить мою подписку на сервис [Название сервиса].

Данные аккаунта:
- Email: [email@example.com]
- [Дополнительные идентификаторы]

Причина отмены: [Выбранная причина / Свой вариант]

Прошу подтвердить отмену подписки и прекращение автоматических списаний.

С уважением,
[Имя пользователя]
```

---

### 📦 МОДУЛЬ 12: Административная панель (Admin Module)

**Описание:** Панель управления для администраторов сервиса.

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| ADM-001 | Аутентификация администраторов | High | 8 | AUTH-001 |
| ADM-002 | Роли и права доступа (RBAC) | High | 16 | ADM-001 |
| ADM-003 | Dashboard со статистикой платформы | High | 16 | - |
| ADM-004 | Управление пользователями | High | 12 | - |
| ADM-005 | Управление каталогом сервисов | High | 12 | CAT-001 |
| ADM-006 | Модерация контента | Medium | 16 | CAT-001 |
| ADM-007 | Логи и аудит действий | High | 16 | - |
| ADM-008 | Управление шаблонами уведомлений | Medium | 12 | NOTIF-001 |
| ADM-009 | Настройки системы | Medium | 8 | - |
| ADM-010 | Управление парсерами почты | Medium | 16 | EMAIL-001 |
| ADM-011 | Мониторинг интеграций | High | 12 | BANK-001, EMAIL-001 |
| ADM-012 | Экспорт аналитики платформы | Low | 10 | ADM-003 |

---

### 📦 МОДУЛЬ 13: Мобильное приложение (Mobile App Module)

**Описание:** Специфичные для мобильных задачи (React Native).

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| MOB-001 | Настройка React Native проекта (Expo / bare workflow) | High | 8 | - |
| MOB-002 | Архитектура (Zustand, Clean Architecture, Feature-Sliced) | High | 16 | MOB-001 |
| MOB-003 | Дизайн-система (компоненты, темы, NativeWind) | High | 24 | MOB-001 |
| MOB-004 | Онбординг новых пользователей | High | 16 | MOB-003 |
| MOB-005 | Главный экран (Dashboard) | High | 20 | MOB-003 |
| MOB-006 | Список подписок | High | 16 | MOB-003, SUB-002 |
| MOB-007 | Детальный экран подписки | High | 12 | MOB-003, SUB-003 |
| MOB-008 | Добавление/редактирование подписки | High | 16 | MOB-003, SUB-001 |
| MOB-009 | Экран аналитики | High | 20 | MOB-003, ANAL-001 |
| MOB-010 | Экран уведомлений | High | 12 | MOB-003, NOTIF-001 |
| MOB-011 | Настройки профиля | Medium | 12 | MOB-003, USER-001 |
| MOB-012 | Подключение банков | High | 16 | MOB-003, BANK-001 |
| MOB-013 | Подключение почты | High | 16 | MOB-003, EMAIL-001 |
| MOB-014 | Push-уведомления | High | 12 | NOTIF-001 |
| MOB-015 | Локальные уведомления | Medium | 8 | - |
| MOB-016 | Биометрическая аутентификация | Medium | 12 | AUTH-003 |
| MOB-017 | Виджеты для домашнего экрана (iOS/Android) | Low | 24 | - |
| MOB-018 | Apple Watch / Wear OS компаньон (React Native Watch Connectivity) | Low | 40 | - |
| MOB-019 | Offline режим | Medium | 20 | - |
| MOB-020 | Deep linking | Medium | 12 | - |
| MOB-021 | App Store / Play Store публикация | High | 16 | - |

---

### 📦 МОДУЛЬ 14: Веб-приложение (Web App Module)

**Описание:** Специфичные для веба задачи (Next.js).

#### Задачи модуля:

| ID | Задача | Приоритет | Оценка (ч) | Зависимости |
|----|--------|-----------|------------|-------------|
| WEB-001 | Настройка Next.js проекта | High | 8 | - |
| WEB-002 | Архитектура (App Router, Server Components) | High | 12 | WEB-001 |
| WEB-003 | Дизайн-система (Tailwind + компоненты) | High | 20 | WEB-001 |
| WEB-004 | Лендинг (главная страница) | High | 24 | WEB-003 |
| WEB-005 | Страницы авторизации | High | 12 | WEB-003, AUTH-001 |
| WEB-006 | Dashboard | High | 20 | WEB-003, ANAL-001 |
| WEB-007 | Список подписок (таблица + карточки) | High | 16 | WEB-003, SUB-002 |
| WEB-008 | Детальная страница подписки | High | 12 | WEB-003, SUB-003 |
| WEB-009 | Форма подписки | High | 12 | WEB-003, SUB-001 |
| WEB-010 | Страница аналитики | High | 24 | WEB-003, ANAL-001 |
| WEB-011 | Страница уведомлений | Medium | 12 | WEB-003, NOTIF-001 |
| WEB-012 | Настройки профиля | Medium | 12 | WEB-003, USER-001 |
| WEB-013 | Подключение интеграций | High | 20 | WEB-003, BANK-001, EMAIL-001 |
| WEB-014 | Каталог сервисов | Medium | 16 | WEB-003, CAT-001 |
| WEB-015 | Сравнение сервисов | Medium | 12 | WEB-003, CAT-005 |
| WEB-016 | PWA (Progressive Web App) | Medium | 16 | - |
| WEB-017 | SEO оптимизация | Medium | 12 | WEB-004 |
| WEB-018 | Accessibility (a11y) | Medium | 16 | WEB-003 |
| WEB-019 | Интернационализация (i18n) | Low | 20 | - |
| WEB-020 | Деплой и CI/CD | High | 16 | - |

---

## 5. План разработки

### 5.1 Фазы проекта

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          TIMELINE ПРОЕКТА                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Фаза 1: MVP                    Фаза 2: Интеграции    Фаза 3: Scale        │
│  ══════════════                 ══════════════════    ═════════════         │
│  │                              │                     │                     │
│  ▼                              ▼                     ▼                     │
│  ┌──────────────────┐          ┌──────────────────┐  ┌──────────────────┐  │
│  │  3 месяца        │          │  2 месяца        │  │  2 месяца        │  │
│  │  ─────────       │          │  ─────────       │  │  ─────────       │  │
│  │  • Auth          │          │  • Банки         │  │  • ML прогнозы   │  │
│  │  • CRUD подписок │          │  • Почта         │  │  • Виджеты       │  │
│  │  • Базовая       │          │  • Авто-создание │  │  • Watch apps    │  │
│  │    аналитика     │          │    подписок      │  │  • Расширенная   │  │
│  │  • Уведомления   │          │  • Альтернативы  │  │    аналитика     │  │
│  │  • Mobile + Web  │          │  • Отмена        │  │  • Партнёрства   │  │
│  └──────────────────┘          └──────────────────┘  └──────────────────┘  │
│           │                             │                     │             │
│           │      Beta Launch            │    Public Launch    │             │
│           └──────────────●──────────────┴──────────●──────────┘             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Фаза 1: MVP (12 недель)

**Цель:** Базовый функционал для ручного управления подписками

#### Спринт 1-2 (недели 1-4)

| Неделя | Backend | Mobile | Web |
|--------|---------|--------|-----|
| 1 | Настройка инфраструктуры, CI/CD | Настройка React Native (Expo) проекта | Настройка Next.js |
| 2 | AUTH-001..005 (регистрация, логин) | MOB-001..003 (архитектура) | WEB-001..003 |
| 3 | AUTH-006..008 (OAuth) | MOB-004 (онбординг) | WEB-004..005 (лендинг, auth) |
| 4 | USER-001..006 (профиль) | MOB-005, MOB-011 | WEB-006, WEB-012 |

#### Спринт 3-4 (недели 5-8)

| Неделя | Backend | Mobile | Web |
|--------|---------|--------|-----|
| 5 | SUB-001..005 (CRUD подписок) | MOB-006..007 | WEB-007..008 |
| 6 | SUB-006..010 (фильтры, категории) | MOB-008 | WEB-009 |
| 7 | ANAL-001..004 (базовая аналитика) | MOB-009 | WEB-010 |
| 8 | PRED-001..002 (базовый прогноз) | MOB-009 (продолжение) | WEB-010 (продолжение) |

#### Спринт 5-6 (недели 9-12)

| Неделя | Backend | Mobile | Web |
|--------|---------|--------|-----|
| 9 | NOTIF-001..004 (уведомления) | MOB-010, MOB-014 | WEB-011 |
| 10 | NOTIF-009..011 (настройки) | MOB-015 | WEB-011 |
| 11 | Тестирование, баг-фиксы | MOB-016, MOB-019 | WEB-016, WEB-17 |
| 12 | Подготовка к релизу | MOB-021 | WEB-020 |

**Deliverables Фазы 1:**
- ✅ Регистрация/авторизация (email + OAuth)
- ✅ CRUD подписок
- ✅ Базовая аналитика и визуализация
- ✅ Напоминания о списаниях
- ✅ Мобильные приложения (iOS + Android)
- ✅ Веб-приложение

---

### 5.3 Фаза 2: Интеграции (8 недель)

**Цель:** Автоматический сбор данных и продвинутые функции

#### Спринт 7-8 (недели 13-16)

| Неделя | Backend | Mobile | Web |
|--------|---------|--------|-----|
| 13 | BANK-001..003 (OAuth банков) | MOB-012 | WEB-013 |
| 14 | BANK-004..006 (транзакции) | MOB-012 | WEB-013 |
| 15 | BANK-007..010 (сопоставление) | Интеграция | Интеграция |
| 16 | EMAIL-001..004 (OAuth почты) | MOB-013 | WEB-013 |

#### Спринт 9-10 (недели 17-20)

| Неделя | Backend | Mobile | Web |
|--------|---------|--------|-----|
| 17 | EMAIL-005..009 (парсинг) | MOB-013 | WEB-013 |
| 18 | EMAIL-010..014 (ML парсинг) | Интеграция | Интеграция |
| 19 | CAT-001..004, ALT-001..002 | WEB-014 | WEB-014 |
| 20 | CANC-001..005 (отмена) | Интеграция | WEB-015 |

**Deliverables Фазы 2:**
- ✅ Интеграция с банками (Сбербанк, Тинькофф, ЮMoney)
- ✅ Парсинг почты (Gmail, Яндекс, Mail.ru)
- ✅ Автоматическое создание подписок
- ✅ Каталог сервисов и альтернативы
- ✅ Помощь в отмене подписок

---

### 5.4 Фаза 3: Масштабирование (8 недель)

**Цель:** Продвинутые функции, оптимизация, расширение

#### Спринт 11-12 (недели 21-24)

| Неделя | Backend | Mobile | Web |
|--------|---------|--------|-----|
| 21 | PRED-003..007 (ML прогнозы) | USE-001..006 | USE-001..006 |
| 22 | ALT-003..006 (рекомендации) | Интеграция | Интеграция |
| 23 | AUTH-009..010 (2FA) | MOB-016 | Интеграция |
| 24 | ANAL-008..012 (экспорты) | MOB-017 | WEB-018 |

#### Спринт 13-14 (недели 25-28)

| Неделя | Backend | Mobile | Web |
|--------|---------|--------|-----|
| 25 | ADM-001..006 | MOB-018 | Админ-панель |
| 26 | ADM-007..012 | MOB-018 | Админ-панель |
| 27 | Оптимизация производительности | Оптимизация | Оптимизация |
| 28 | Подготовка к публичному релизу | Финальное тестирование | Финальное тестирование |

**Deliverables Фазы 3:**
- ✅ ML-прогнозирование затрат
- ✅ Отслеживание использования
- ✅ 2FA аутентификация
- ✅ Виджеты для домашнего экрана
- ✅ Административная панель
- ✅ Оптимизация производительности

---

### 5.5 Команда

| Роль | Количество | Ответственность |
|------|------------|-----------------|
| Project Manager | 1 | Управление проектом, сроки, коммуникация |
| Tech Lead / Architect | 1 | Архитектура, code review, техническое лидерство |
| Backend Developer | 2 | Go, API, интеграции, БД |
| Mobile Developer (React Native) | 2 | iOS/Android приложения |
| Frontend Developer | 2 | Next.js веб-приложение |
| DevOps Engineer | 1 | Инфраструктура, CI/CD, мониторинг |
| UI/UX Designer | 1 | Дизайн, прототипы, дизайн-система |
| QA Engineer | 1 | Тестирование, автотесты |
| ML Engineer | 0.5 | Парсинг почты, прогнозирование (part-time/фаза 2-3) |

**Итого:** 11.5 человек

---

## 6. Требования к безопасности

### 6.1 Аутентификация и авторизация

- [ ] JWT токены с коротким временем жизни (15 мин access, 7 дней refresh)
- [ ] Secure, HttpOnly, SameSite cookies
- [ ] Rate limiting на auth endpoints
- [ ] Защита от brute-force (блокировка после 5 попыток)
- [ ] PKCE для OAuth на мобильных

### 6.2 Хранение данных

- [ ] Шифрование sensitive данных в БД (AES-256)
- [ ] Отдельное хранилище для секретов (HashiCorp Vault)
- [ ] Регулярные бекапы БД (ежедневно)
- [ ] Ротация ключей шифрования

### 6.3 Сетевая безопасность

- [ ] HTTPS everywhere (TLS 1.3)
- [ ] CORS политика
- [ ] CSP, X-Frame-Options, X-Content-Type-Options
- [ ] API rate limiting
- [ ] WAF (Web Application Firewall)

### 6.4 Конфиденциальность

- [ ] Минимизация сбора данных
- [ ] Явное согласие на доступ к почте/банкам
- [ ] Возможность удаления всех данных (GDPR)
- [ ] Политика конфиденциальности
- [ ] Соответствие 152-ФЗ (персональные данные)

### 6.5 Аудит и мониторинг

- [ ] Логирование всех действий с чувствительными данными
- [ ] Мониторинг аномалий
- [ ] Регулярные security-аудиты
- [ ] Penetration testing перед релизом

---

## 7. Нефункциональные требования

### 7.1 Производительность

| Метрика | Целевое значение |
|---------|------------------|
| API response time (p95) | < 200ms |
| Время загрузки мобильного приложения | < 2s |
| Время загрузки веб-страницы (LCP) | < 2.5s |
| API availability | 99.9% |
| Concurrent users | 10,000+ |

### 7.2 Масштабируемость

- Горизонтальное масштабирование микросервисов
- Auto-scaling на основе нагрузки
- CDN для статических ресурсов
- Кэширование на всех уровнях

### 7.3 Совместимость

| Платформа | Требования |
|-----------|------------|
| iOS | 14.0+, iPhone 8+, iPad Air 3+ |
| Android | 8.0+, ARM64, 2GB RAM+ |
| Web | Chrome 90+, Firefox 88+, Safari 14+, Edge 90+ |

### 7.4 Локализация

- Русский язык (основной)
- Английский язык (Phase 3)
- Поддержка валют: RUB, USD, EUR, KZT, BYN

### 7.5 Accessibility

- WCAG 2.1 Level AA
- Поддержка screen readers
- Альтернативный текст для изображений
- Достаточный контраст цветов
- Навигация с клавиатуры (веб)

---

## 📊 Сводная таблица модулей

| # | Модуль | Задач | Часов | Приоритет |
|---|--------|-------|-------|-----------|
| 1 | Auth | 12 | 112 | Critical |
| 2 | User Profile | 9 | 62 | High |
| 3 | Subscriptions | 14 | 106 | Critical |
| 4 | Bank Integration | 13 | 178 | High |
| 5 | Email Parsing | 16 | 252 | High |
| 6 | Analytics | 12 | 130 | High |
| 7 | Usage Tracking | 10 | 114 | Medium |
| 8 | Prediction | 10 | 106 | Medium |
| 9 | Notifications | 15 | 150 | High |
| 10 | Catalog & Alternatives | 11 | 168 | Medium |
| 11 | Cancellation | 10 | 114 | Medium |
| 12 | Admin | 12 | 136 | Medium |
| 13 | Mobile App | 21 | 296 | Critical |
| 14 | Web App | 20 | 268 | Critical |
| **Итого** | | **185** | **2,192** | |

**Оценка с учётом коэффициента 1.5x на риски и коммуникацию:**
- Общая оценка: ~3,288 часов
- При команде 10 разработчиков: ~7 месяцев

---

## 📎 Приложения

### Приложение A: User Stories

<details>
<summary>Развернуть User Stories</summary>

**Регистрация и вход:**
- Как пользователь, я хочу зарегистрироваться через email, чтобы создать аккаунт
- Как пользователь, я хочу войти через Google/Яндекс, чтобы упростить авторизацию
- Как пользователь, я хочу восстановить пароль, если забыл его

**Подписки:**
- Как пользователь, я хочу добавить подписку вручную, чтобы отслеживать расходы
- Как пользователь, я хочу видеть все подписки в одном месте
- Как пользователь, я хочу фильтровать подписки по категориям
- Как пользователь, я хочу видеть дату следующего списания

**Интеграции:**
- Как пользователь, я хочу подключить банк, чтобы подписки добавлялись автоматически
- Как пользователь, я хочу подключить почту для парсинга чеков
- Как пользователь, я хочу подтверждать найденные подписки перед добавлением

**Аналитика:**
- Как пользователь, я хочу видеть общую сумму затрат на подписки
- Как пользователь, я хочу видеть распределение по категориям
- Как пользователь, я хочу видеть прогноз затрат на год

**Уведомления:**
- Как пользователь, я хочу получать напоминание за 3 дня до списания
- Как пользователь, я хочу узнавать об изменении цены
- Как пользователь, я хочу настроить типы уведомлений

</details>

### Приложение B: Глоссарий

| Термин | Определение |
|--------|-------------|
| Подписка | Регулярный платёж за доступ к сервису |
| Биллинговый цикл | Периодичность списания (месяц, год, неделя) |
| Парсинг | Извлечение структурированных данных из текста |
| OAuth | Протокол авторизации для доступа к сторонним сервисам |
| JWT | JSON Web Token - формат токена авторизации |
| Push-уведомление | Уведомление на устройство пользователя |

---

**Документ подготовлен:** MultiSub Team  
**Контакт:** team@multisub.app  
**Версия:** 1.1 (MVP)

---

## 🧠 AI Agent Skills для проекта

> Скиллы подключены в `.agents/skills/` и автоматически используются GitHub Copilot, Cursor, Codex и другими AI-агентами в этом проекте.  
> Если нужно добавить новый скилл: `npx --yes skills add <owner/repo@skill-name> -y`

| # | Скилл | Цель |
|---|-------|------|
| 1 | `vercel-react-best-practices` | Лучшие практики React / Next.js от Vercel Engineering |
| 2 | `vercel-react-native-skills` | React Native: паттерны, перформанс, типичные ошибки |
| 3 | `vercel-composition-patterns` | Server / Client Component паттерны Next.js App Router |

**Установка дополнительных скиллов:**

```bash
# React Native от Callstack
npx --yes skills add callstackincubator/agent-skills@react-native-best-practices -y

# Next.js App Router паттерны
npx --yes skills add wshobson/agents@nextjs-app-router-patterns -y

# GitHub Actions шаблоны
npx --yes skills add wshobson/agents@github-actions-templates -y
```
