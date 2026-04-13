-- +goose Up
-- +goose StatementBegin
CREATE TABLE problems_tests (
    id              BIGSERIAL PRIMARY KEY,
    problem_id      UUID REFERENCES problems(id),
    input_data      TEXT NOT NULL,
    output_data     TEXT NOT NULL,
    is_hidden       BOOLEAN NOT NULL DEFAULT false,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX problems_tests_problem_id_idx ON problems_tests(problem_id);

COMMENT ON TABLE problems_tests IS 'Тесты к задачам';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE problems_tests;
-- +goose StatementEnd
