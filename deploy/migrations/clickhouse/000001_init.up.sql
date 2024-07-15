CREATE TABLE IF NOT EXISTS history_orders (
    client_name String,
    exchange_name String,
    label String,
    pair String,
    side String,
    type String,
    base_qty Float64,
    price Float64,
    algorithm_name_placed String,
    lowest_sell_prc Float64,
    highest_buy_prc Float64,
    commission_quote_qty Float64,
    time_placed DateTime
)
ENGINE = MergeTree()
ORDER BY (exchange_name, pair, label, client_name);