-- name: CreateOrUpdateAsset :one
insert into
    assets (code, price, updated_at)
values
    ($1, $2, $3)
on conflict (code)
    do update set
                  price      = excluded.price,
                  updated_at = excluded.updated_at
returning *;

-- name: CreateAssetPriceLog :one
insert into
    asset_price_logs (asset_id, price)
values
    ($1, $2)
returning *;

-- name: ListAssetsWithPriceChanges :many
with
    price_changes_4h as (select distinct on (asset_id)
                             asset_id,
                             first_value(price)
                             over (partition by asset_id order by created_at range between unbounded preceding and unbounded following) as first,
                             last_value(price)
                             over (partition by asset_id order by created_at range between unbounded preceding and unbounded following) as last
                         from
                             asset_price_logs
                         where
                             created_at > now() - interval '4 hours'),
    price_changes_1d as (select distinct on (asset_id)
                             asset_id,
                             first_value(price)
                             over (partition by asset_id order by created_at range between unbounded preceding and unbounded following) as first,
                             last_value(price)
                             over (partition by asset_id order by created_at range between unbounded preceding and unbounded following) as last
                         from
                             asset_price_logs
                         where
                             created_at > now() - interval '1 day'),
    price_changes_7d as (select distinct on (asset_id)
                             asset_id,
                             first_value(price)
                             over (partition by asset_id order by created_at range between unbounded preceding and unbounded following) as first,
                             last_value(price)
                             over (partition by asset_id order by created_at range between unbounded preceding and unbounded following) as last
                         from
                             asset_price_logs
                         where
                             created_at > now() - interval '7 days')

select *
from
    assets
        join lateral (
        select
            price_changes_4h.first::numeric                           as first4h,
            (price_changes_4h.last - price_changes_4h.first)::numeric as change4h
        from
            price_changes_4h
        where
            asset_id = assets.id
        ) as pch4h on true

        join lateral (
        select
            price_changes_1d.first::numeric                           as first1d,
            (price_changes_1d.last - price_changes_1d.first)::numeric as change1d
        from
            price_changes_1d
        where
            asset_id = assets.id
        ) as pch1d on true

        join lateral (
        select
            price_changes_7d.first::numeric                           as first7d,
            (price_changes_7d.last - price_changes_7d.first)::numeric as change7d
        from
            price_changes_7d
        where
            asset_id = assets.id
        ) as pch7d on true;

-- name: GetAsset :one
select *
from
    assets
where
    code = $1;