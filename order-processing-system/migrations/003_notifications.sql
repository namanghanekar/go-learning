create table if not exists notifications (
    id bigserial primary key,
    order_id text not null,
    event_id text not null unique,
    event_type text not null,
    payload jsonb,
    created_at timestamptz not null default now()
);

create index if not exists idx_notifications_order_id on notifications(order_id);
