-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists tasks (
    id bigserial primary key,
    task_type bigint not null default 0,
    payload JSONB not null,
    status bigint not null default 0,
    created_at timestamptz not null default NOW(),
    updated_at timestamptz default null,
    completed_at timestamptz default null,
    attempts_left bigint default 3
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists tasks;
-- +goose StatementEnd
