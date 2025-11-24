-- name: CreateOrUpdateAsset :one
insert into
    assets (code, price, updated_at)
values
    ($1, $2, $3)
on conflict (code) do update set
    price = excluded.price
returning *;

-- name: CreateAssetPriceLog :one
insert into
    asset_price_logs (asset_id, price)
values
    ($1, $2)
returning *;

-- name: ListAssets :many
select *
from
    assets;