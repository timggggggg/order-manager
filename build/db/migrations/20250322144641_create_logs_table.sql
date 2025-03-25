-- +goose Up
-- +goose StatementBegin
create table if not exists logs (
    id bigserial primary key,
    order_id bigint not null,
    status_from ostatus not null,
    status_to ostatus not null,
    ts timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists logs;
-- +goose StatementEnd
