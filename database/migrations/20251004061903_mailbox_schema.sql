-- +goose Up
-- +goose StatementBegin

CREATE TYPE mailbox_status AS ENUM ('active', 'inactive', 'suspended');
CREATE TYPE contact_status AS ENUM ('active', 'unsubscribed');
CREATE TYPE sequence_contact_status AS ENUM ('pending', 'in_progress', 'completed', 'paused', 'bounced', 'cancelled');
CREATE TYPE email_queue_status AS ENUM ('scheduled', 'queued', 'sending', 'sent', 'failed', 'cancelled');
CREATE TYPE email_event_type AS ENUM ('sent', 'delivered', 'opened', 'clicked', 'bounced', 'failed');

CREATE TABLE mailboxes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    daily_capacity INTEGER DEFAULT 30,
    status mailbox_status DEFAULT 'active',
    provider VARCHAR(100),
    smtp_host VARCHAR(255),
    smtp_port INTEGER,
    smtp_username VARCHAR(255),
    encrypted_smtp_password BYTEA,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE sequence_mailboxes (
    sequence_id UUID REFERENCES sequences(id) ON DELETE CASCADE,
    mailbox_id UUID REFERENCES mailboxes(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (sequence_id, mailbox_id)
);

CREATE TABLE contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    company VARCHAR(255),
    phone VARCHAR(50),
    status contact_status DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE sequence_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sequence_id UUID NOT NULL REFERENCES sequences(id) ON DELETE CASCADE,
    contact_id UUID NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    current_step INTEGER DEFAULT 0,
    next_send_at TIMESTAMP WITH TIME ZONE,
    status sequence_contact_status DEFAULT 'pending',
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(sequence_id, contact_id)
);

CREATE TABLE email_queues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sequence_contact_id UUID NOT NULL REFERENCES sequence_contacts(id) ON DELETE CASCADE,
    mailbox_id UUID REFERENCES mailboxes(id) ON DELETE SET NULL,
    step_order INTEGER NOT NULL,
    subject TEXT NOT NULL,
    content TEXT NOT NULL,
    scheduled_for TIMESTAMP WITH TIME ZONE NOT NULL,
    status email_queue_status DEFAULT 'scheduled',
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    last_attempt_at TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE mailbox_daily_counts (
    mailbox_id UUID REFERENCES mailboxes(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    sent_count INTEGER DEFAULT 0,
    failed_count INTEGER DEFAULT 0,
    reset_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (mailbox_id, date)
);

CREATE TABLE email_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email_queue_id UUID REFERENCES email_queues(id) ON DELETE CASCADE,
    event_type email_event_type NOT NULL,
    event_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE kafka_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic VARCHAR(100) NOT NULL,
    partition INTEGER NOT NULL,
    offset BIGINT NOT NULL,
    processed_at TIMESTAMPTZ DEFAULT NOW(),
    email_count INTEGER DEFAULT 0,
    batch_size INTEGER DEFAULT 0
);

CREATE INDEX idx_sequence_mailboxes_mailbox_id ON sequence_mailboxes(mailbox_id);
CREATE INDEX idx_contacts_email ON contacts(email);
CREATE INDEX idx_contacts_status ON contacts(status);
CREATE INDEX idx_sequence_contacts_sequence_id ON sequence_contacts(sequence_id);
CREATE INDEX idx_sequence_contacts_contact_id ON sequence_contacts(contact_id);
CREATE INDEX idx_sequence_contacts_status ON sequence_contacts(status);
CREATE INDEX idx_sequence_contacts_next_send ON sequence_contacts(next_send_at) WHERE status = 'in_progress';
CREATE INDEX idx_email_queues_scheduled_status ON email_queues(scheduled_for, status);
CREATE INDEX idx_email_queues_sequence_contact ON email_queues(sequence_contact_id);
CREATE INDEX idx_email_queues_mailbox_date ON email_queues(mailbox_id, scheduled_for);
CREATE INDEX idx_mailbox_daily_counts_date ON mailbox_daily_counts(date);
CREATE INDEX idx_email_events_email_queue ON email_events(email_queue_id);
CREATE INDEX idx_email_events_type ON email_events(event_type);

-- Create updated_at triggers
CREATE TRIGGER update_mailboxes_updated_at
    BEFORE UPDATE ON mailboxes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_contacts_updated_at
    BEFORE UPDATE ON contacts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sequence_contacts_updated_at
    BEFORE UPDATE ON sequence_contacts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_email_queues_updated_at
    BEFORE UPDATE ON email_queues
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop triggers first
DROP TRIGGER IF EXISTS update_email_queues_updated_at ON email_queues;
DROP TRIGGER IF EXISTS update_sequence_contacts_updated_at ON sequence_contacts;
DROP TRIGGER IF EXISTS update_contacts_updated_at ON contacts;
DROP TRIGGER IF EXISTS update_mailboxes_updated_at ON mailboxes;

DROP INDEX IF EXISTS idx_email_events_type;
DROP INDEX IF EXISTS idx_email_events_email_queue;
DROP INDEX IF EXISTS idx_mailbox_daily_counts_date;
DROP INDEX IF EXISTS idx_email_queues_mailbox_date;
DROP INDEX IF EXISTS idx_email_queues_sequence_contact;
DROP INDEX IF EXISTS idx_email_queues_scheduled_status;
DROP INDEX IF EXISTS idx_sequence_contacts_next_send;
DROP INDEX IF EXISTS idx_sequence_contacts_status;
DROP INDEX IF EXISTS idx_sequence_contacts_contact_id;
DROP INDEX IF EXISTS idx_sequence_contacts_sequence_id;
DROP INDEX IF EXISTS idx_contacts_status;
DROP INDEX IF EXISTS idx_contacts_email;
DROP INDEX IF EXISTS idx_sequence_mailboxes_mailbox_id;

DROP TABLE IF EXISTS email_events;
DROP TABLE IF EXISTS mailbox_daily_counts;
DROP TABLE IF EXISTS email_queues;
DROP TABLE IF EXISTS sequence_contacts;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS sequence_mailboxes;
DROP TABLE IF EXISTS mailboxes;
DROP TABLE IF EXISTS kafka_batches;

DROP TYPE IF EXISTS email_event_type;
DROP TYPE IF EXISTS email_queue_status;
DROP TYPE IF EXISTS sequence_contact_status;
DROP TYPE IF EXISTS contact_status;
DROP TYPE IF EXISTS mailbox_status;
-- +goose StatementEnd
