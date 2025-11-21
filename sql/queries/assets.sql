-- name: CreateAsset :one
insert into
    assets (code, name, price)
values
    ($1, $2, $3)
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