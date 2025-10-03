-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE sequences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    open_tracking_enabled BOOLEAN DEFAULT true,
    click_tracking_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sequence_id UUID NOT NULL REFERENCES sequences(id) ON DELETE CASCADE,
    step_order INTEGER NOT NULL,
    subject TEXT NOT NULL,
    content TEXT NOT NULL,
    wait_days INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_steps_sequence_id ON steps(sequence_id);
CREATE INDEX idx_steps_order ON steps(sequence_id, step_order);
CREATE INDEX idx_sequences_deleted_at ON sequences(deleted_at);
CREATE INDEX idx_steps_deleted_at ON steps(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_steps_order;
DROP INDEX IF EXISTS idx_steps_sequence_id;
DROP TABLE IF EXISTS steps;
DROP TABLE IF EXISTS sequences;
-- +goose StatementEnd
