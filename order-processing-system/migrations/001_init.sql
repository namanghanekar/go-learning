create table if not exists orders (
    id text primary key,
    user_id text not null,
    status text not null,
    amount_cents bigint not null,
    items jsonb not null,
    email text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_orders_user_id on orders(user_id);
create index if not exists idx_orders_status on orders(status);
