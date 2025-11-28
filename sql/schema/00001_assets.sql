-- +goose Up
-- +goose StatementBegin
create table assets
(
    id         bigint generated always as identity primary key,
    code       varchar(20)    not null unique,
    price      numeric(24, 6) not null,
    created_at timestamptz    not null default now(),
    updated_at timestamptz    not null
);

create index assets_created_at_idx on assets (created_at);

create table asset_price_logs
(
    id         bigint generated always as identity primary key,
    asset_id   bigint references assets (id),
    price      numeric(24, 6) not null,
    created_at timestamptz    not null default now()
);

create index asset_price_logs_created_at_idx on asset_price_logs (created_at);
create index asset_price_logs_asset_id_idx on asset_price_logs (asset_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop table asset_price_logs;

drop table assets;
-- +goose StatementEnd
