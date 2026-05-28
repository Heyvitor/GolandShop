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
