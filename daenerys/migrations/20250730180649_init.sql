-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE submissions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    problem_id  UUID NOT NULL,
    user_id     UUID NOT NULL,
    status      INT NOT NULL DEFAULT 0,
    code        TEXT NOT NULL,
    language    INT NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX submissions_problem_id_user_id_idx ON submissions(problem_id, user_id);

COMMENT ON TABLE submissions IS 'Отправленные решения пользователя';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE submissions;
-- +goose StatementEnd
