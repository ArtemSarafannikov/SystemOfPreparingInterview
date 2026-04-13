-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE problems (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_id       UUID,
    summary         TEXT NOT NULL,
    description     text,
    time_limit_ms   INTEGER NOT NULL,
    memory_limit_kb INTEGER NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX problems_author_id_idx ON problems(author_id);

COMMENT ON TABLE problems IS 'Хранилище задач';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE problems;
-- +goose StatementEnd
