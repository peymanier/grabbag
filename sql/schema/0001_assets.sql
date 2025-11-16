-- +goose Up
create table assets
(
    id         bigint generated always as identity primary key,
    code       varchar(5)     not null unique,
    name       varchar(50)    not null,
    price      numeric(15, 6) not null,
    created_at timestamptz    not null default now()
);

create table asset_price_logs
(
    id         bigint generated always as identity primary key,
    asset_id   bigint references assets (id),
    price      numeric(15, 6) not null,
    created_at timestamptz    not null default now()
);

create index asset_values_from_asset_idx on asset_price_logs (asset_id);

-- +goose Down
drop table assets;

drop table asset_price_logs;
