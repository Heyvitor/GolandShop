-- =============================================================================
-- SQL DE MIGRAÇÃO MANUAL (RODAR NO TERMINAL DO POSTGRES)
-- =============================================================================

-- 1. Garante que a extensão de UUID existe
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 2. Adiciona a coluna de Role que está causando o erro de scan no Go
ALTER TABLE users ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'client';

-- 3. Cria a tabela de Lojas (Multi-Tenancy)
CREATE TABLE IF NOT EXISTS stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 4. Cria a tabela de Pedidos (Orders)
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 5. Criação de índices para performance em alta escala
CREATE INDEX IF NOT EXISTS idx_stores_owner_id ON stores(owner_id);
CREATE INDEX IF NOT EXISTS idx_stores_slug ON stores(slug);
CREATE INDEX IF NOT EXISTS idx_orders_client_id ON orders(client_id);
CREATE INDEX IF NOT EXISTS idx_orders_store_id ON orders(store_id);

-- 6. Evolui a tabela de items para produtos completos
ALTER TABLE items
    ADD COLUMN IF NOT EXISTS store_id UUID REFERENCES stores(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS name TEXT,
    ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    ADD COLUMN IF NOT EXISTS variant TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS variant_price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    ADD COLUMN IF NOT EXISTS shipping_type TEXT NOT NULL DEFAULT 'consult';

UPDATE items
SET
    name = COALESCE(NULLIF(name, ''), title),
    description = COALESCE(description, body, ''),
    store_id = COALESCE(store_id, (
        SELECT s.id
        FROM stores s
        WHERE s.owner_id = items.user_id
        ORDER BY s.created_at ASC
        LIMIT 1
    ))
WHERE name IS NULL OR store_id IS NULL;

ALTER TABLE items
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN store_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_items_store_id_created_at ON items(store_id, created_at DESC);

-- Opcional: Se você já tem um usuário administrador e quer mudar a role dele manualmente:
-- UPDATE users SET role = 'admin' WHERE email = 'seu-email@exemplo.com';
