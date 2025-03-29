-- +goose Up
-- +goose StatementBegin
create table if not exists http_req_logs (
    request_id bigserial primary key,
    ts timestamptz,
    method varchar(10) not null,
    url text,
    request_body text
);

create table if not exists http_resp_logs (
    response_id bigserial primary key,
    ts timestamptz,
    status_code int not null check (status_code between 100 and 599),
    body text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists http_resp_logs;
drop table if exists http_req_logs;
-- +goose StatementEnd
