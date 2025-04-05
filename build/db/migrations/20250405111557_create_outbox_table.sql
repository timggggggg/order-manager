-- +goose Up
-- +goose StatementBegin
create type task_status as enum (
    'CREATED',
    'PROCESSING',
    'COMPLETED',
    'FAILED',
    'NO_ATTEMPTS_LEFT'
);

CREATE TABLE if not exists outbox (
    id bigserial primary key,
    audit_log JSONB not null,
    status task_status not null default 'CREATED',
    created_at timestamptz not null default NOW(),
    updated_at timestamptz default null,
    completed_at timestamptz default null,
    attempts_left bigint default 3
);

DO $$
BEGIN
    INSERT INTO outbox (audit_log, status, created_at)
    SELECT 
        jsonb_build_object(
            'order_id', a.order_id,
            'status_from', a.status_from,
            'status_to', a.status_to,
            'ts', a.ts
        ),
        'CREATED',
        NOW()
    FROM logs a;

    -- ts timestamptz,
    -- method varchar(10) not null,
    -- url text,
    -- request_body text
    INSERT INTO outbox (audit_log, status, created_at)
    SELECT 
        jsonb_build_object(
            'ts', a.ts,
            'method', a.method,
            'url', a.url,
            'request_body', a.request_body
        ),
        'CREATED',
        NOW()
    FROM http_req_logs a;

    -- ts timestamptz,
    -- status_code int not null check (status_code between 100 and 599),
    -- body text
    INSERT INTO outbox (audit_log, status, created_at)
    SELECT 
        jsonb_build_object(
            'ts', a.ts,
            'status_code', a.status_code,
            'body', a.body
        ),
        'CREATED',
        NOW()
    FROM http_resp_logs a;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists outbox;
drop type if exists task_status;
-- +goose StatementEnd
