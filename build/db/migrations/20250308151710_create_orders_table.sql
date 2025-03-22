-- +goose Up
-- +goose StatementBegin
create type packaging_type as enum (
    '',
    'bag',
    'box',
    'film'
);

create type ostatus as enum (
    'nil',
    'accepted',
    'expired',
    'issued',
    'returned',
    'withdrawed'
);

create table if not exists orders (
    id bigserial primary key,
    user_id bigint not null,
    order_status ostatus not null default 'accepted',
    accept_time timestamptz,
    expire_time timestamptz,
    issue_time timestamptz,

    weight float not null,
    cost varchar(255) not null,
    package packaging_type not null default '',
    extra_package packaging_type
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists orders;
drop type if exists ostatus;
drop type if exists packaging_type;
-- +goose StatementEnd
