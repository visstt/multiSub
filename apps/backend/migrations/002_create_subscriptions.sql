-- +goose Up

-- Категории
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    icon VARCHAR(50),
    parent_id UUID REFERENCES categories(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Сервисы (каталог)
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
    is_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_services_category ON services(category_id);

-- Подписки пользователей
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_id UUID REFERENCES services(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
    billing_cycle VARCHAR(20) NOT NULL, -- monthly, yearly, weekly, quarterly
    next_billing_date DATE,
    start_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, paused, cancelled
    auto_renew BOOLEAN NOT NULL DEFAULT true,
    source VARCHAR(50), -- manual, email_parse, bank_sync
    category_id UUID REFERENCES categories(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(user_id, status);
CREATE INDEX idx_subscriptions_next_billing ON subscriptions(next_billing_date) WHERE status = 'active';

-- История платежей
CREATE TABLE payment_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
    payment_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL, -- success, failed, pending
    source VARCHAR(50), -- bank_sync, email_parse, manual
    raw_data JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payment_history_sub ON payment_history(subscription_id);
CREATE INDEX idx_payment_history_date ON payment_history(payment_date);

-- +goose Down
DROP TABLE IF EXISTS payment_history;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS categories;
