DO $$ BEGIN
    CREATE TYPE depth_order AS (price FLOAT8, base_qty FLOAT8);
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS order_books (
    id SERIAL PRIMARY KEY,
    exchange TEXT NOT NULL,
    pair TEXT NOT NULL,
    bids depth_order[] NOT NULL,
    asks depth_order[] NOT NULL,
    UNIQUE(exchange, pair)
);
CREATE INDEX IF NOT EXISTS order_books_exchange_pair_idx ON order_books (exchange, pair);
