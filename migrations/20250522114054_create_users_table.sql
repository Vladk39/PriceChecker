-- +goose Up
-- +goose StatementBegin
CREATE TABLE if NOT exists public.currency(
    symbol VARCHAR(10) PRIMARY KEY,
    price_24h decimal,
    volume_24h decimal,
    last_trade_price decimal
);

INSERT INTO public.currency (symbol, price_24h, volume_24h, last_trade_price) VALUES 
('BTC-USD', 0, 0, 0),
('ETH-USD', 0, 0, 0),
('BTC-TRY', 0, 0, 0),
('ETH-TRY', 0, 0, 0),
('BTC-GBP', 0, 0, 0),
('ETH-GBP', 0, 0, 0),
('BTC-EUR', 0, 0, 0),
('ETH-EUR', 0, 0, 0),
('BTC-USDT', 0, 0, 0),
('ETH-USDT', 0, 0, 0)
ON CONFLICT (symbol) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd
