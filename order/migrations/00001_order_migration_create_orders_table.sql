-- +goose Up

create table orders(
    uuid UUID primary key default gen_random_uuid(),
    user_uuid UUID not null,
    part_uuid UUID[] not null,
    total_price numeric(10, 2),
    transaction_uuid UUID null,
    payment_method varchar null,
    status varchar not null,
    created_at timestamp default now(),
    updated_at timestamp default now()
    );

-- +goose Down
drop table if exists orders;
