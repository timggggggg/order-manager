-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    INSERT INTO tasks (payload, status, created_at)
    SELECT 
        jsonb_build_object(
            'order_id', a.order_id,
            'status_from', a.status_from,
            'status_to', a.status_to,
            'ts', a.ts
        ),
        0,
        NOW()
    FROM logs a;

    -- ts timestamptz,
    -- method varchar(10) not null,
    -- url text,
    -- request_body text
    INSERT INTO tasks (payload, status, created_at)
    SELECT 
        jsonb_build_object(
            'ts', a.ts,
            'method', a.method,
            'url', a.url,
            'request_body', a.request_body
        ),
        0,
        NOW()
    FROM http_req_logs a;

    -- ts timestamptz,
    -- status_code int not null check (status_code between 100 and 599),
    -- body text
    INSERT INTO tasks (payload, status, created_at)
    SELECT 
        jsonb_build_object(
            'ts', a.ts,
            'status_code', a.status_code,
            'body', a.body
        ),
        0,
        NOW()
    FROM http_resp_logs a;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
