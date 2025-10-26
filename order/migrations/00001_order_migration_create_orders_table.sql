-- +goose Up
create type order_status as ENUM ('PENDING_PAYMENT', 'PAID', 'CANCELLED', 'COMPLETED');
create type payment_method as ENUM ('UNKNOWN', 'CARD', 'SBP', 'CREDIT_CARD', 'INVESTOR_MONEY');

create table orders(
    uuid UUID primary key default gen_random_uuid(),
    user_uuid UUID not null,
    part_uuid UUID[] not null,
    total_price numeric(10, 2),
    transaction_uuid UUID null,
    payment_method payment_method null,
    status order_status not null,
    created_at timestamp default now(),
    updated_at timestamp default now()
    );

-- +goose Down
drop table if exists orders;
