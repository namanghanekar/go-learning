alter table if exists orders
    add column if not exists recipient_name text,
    add column if not exists shipping_address text,
    add column if not exists payment_id text,
    add column if not exists payment_status text,
    add column if not exists shipment_id text,
    add column if not exists item_count integer not null default 0;
