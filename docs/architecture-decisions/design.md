# Email Sequence System

## System Flow Overview

Email sequencing system with Kafka for async email processing. Cron jobs push email batches to Kafka topics. and the consumer process the mail daily 

## 1. Sequence Creation

```mermaid
sequenceDiagram
    title Create Sequence

    participant Client as Client
    participant API as Sequence API
    participant DB as Database

    Client->>API: POST /api/v1/sequences
    Note over Client,API: { name, steps[], open_tracking_enabled, click_tracking_enabled }

    API->>DB: Validate request data
    API->>DB: BEGIN TRANSACTION
    API->>DB: INSERT INTO sequences (name, open_tracking_enabled, click_tracking_enabled)
    loop For each step
        API->>DB: INSERT INTO steps (sequence_id, step_order, subject, content, wait_days)
    end
    API->>DB: COMMIT TRANSACTION

    DB-->>API: Sequence Created with ID
    API-->>Client: 201 Created<br/>{ id, name, step_count }
```

## 2. Add Contacts to Sequence

```mermaid
sequenceDiagram
    title Add Contacts to Sequence

    participant Client as Client
    participant API as Sequence API
    participant DB as Database

    Client->>API: POST /api/v1/sequences/{id}/contacts
    Note over Client,API: { contacts: [{email, first_name, last_name, company}] }

    API->>DB: Validate sequence exists and is active
    DB-->>API: Sequence details with steps

    API->>DB: BEGIN TRANSACTION
    loop For each contact
        API->>DB: INSERT INTO sequence_contacts (sequence_id, contact_id, current_step, status)
        API->>DB: INSERT INTO email_queues (sequence_contact_id, step_order, scheduled_for, status)
    end
    API->>DB: COMMIT TRANSACTION

    DB-->>API: Contacts added successfully
    API-->>Client: 201 Created<br/>{ contacts_added: count }
```

## 3. Cron Job - Email Batch Producer

```mermaid
sequenceDiagram
    title Cron Job - Email Batch Producer

    participant Cron as Cron Job
    participant DB as Database
    participant Kafka as Kafka Producer

    Note over Cron: Runs every 5 minutes

    Cron->>DB: SELECT eq.*, sc.contact_id, c.email, c.first_name,<br/>s.subject, s.content, m.smtp_config<br/>FROM email_queues eq<br/>JOIN sequence_contacts sc ON eq.sequence_contact_id = sc.id<br/>JOIN contacts c ON sc.contact_id = c.id<br/>JOIN steps s ON eq.step_order = s.step_order<br/>JOIN mailboxes m ON eq.mailbox_id = m.id<br/>WHERE eq.scheduled_for <= NOW()<br/>AND eq.status = 'scheduled'<br/>AND m.status = 'active'
    DB-->>Cron: List of ready emails with all data

    Cron->>DB: SELECT m.id, m.daily_capacity,<br/>COALESCE(mdc.sent_count, 0) as sent_today<br/>FROM mailboxes m<br/>LEFT JOIN mailbox_daily_counts mdc ON m.id = mdc.mailbox_id<br/>AND mdc.date = CURRENT_DATE<br/>WHERE m.status = 'active'
    DB-->>Cron: Mailbox capacities

    Cron->>Cron: Filter emails by available capacity
    Note over Cron: sent_today < daily_capacity

    loop For each batch of 50 emails
        Cron->>Kafka: Produce to email-jobs topic
        Note over Kafka: { batch_id, emails: [email_data], produced_at }
    end

    Cron->>DB: UPDATE email_queues<br/>SET status = 'queued'<br/>WHERE id IN (processed_email_ids)
```

## 4. Kafka Consumer - Email Processor

```mermaid
sequenceDiagram
    title Kafka Consumer - Email Processor

    participant Kafka as Kafka
    participant Consumer as Email Consumer
    participant DB as Database
    participant SMTP as SMTP Server

    Note over Consumer: Consumer group

    Kafka->>Consumer: Consume from email-jobs topic
    Note over Kafka: Batch of emails with template data

    loop For each email in batch
        Consumer->>DB: UPDATE email_queues<br/>SET status = 'sending',<br/>last_attempt_at = NOW()<br/>WHERE id = email_queue_id

        Consumer->>SMTP: Connect to mailbox SMTP
        Consumer->>SMTP: Send email with template variables
        Note over SMTP: Subject: {subject}<br/>Body: {content} with tracking

        alt SMTP Success (250 OK)
            SMTP-->>Consumer: Email delivered
            Consumer->>DB: UPDATE email_queues<br/>SET status = 'sent',<br/>sent_at = NOW()<br/>WHERE id = email_queue_id

            Consumer->>DB: INSERT INTO mailbox_daily_counts<br/>(mailbox_id, date, sent_count)<br/>VALUES (mailbox_id, TODAY, 1)<br/>ON CONFLICT UPDATE sent_count = sent_count + 1

            Consumer->>DB: INSERT INTO email_events<br/>(email_queue_id, event_type, event_data)<br/>VALUES (email_queue_id, 'sent', '{"timestamp": NOW()}')

            Consumer->>Kafka: Produce to followup-events topic
            Note over Kafka: { sequence_contact_id: id, step_completed: step_order }

        else SMTP Failure
            SMTP-->>Consumer: Error
            Consumer->>DB: UPDATE email_queues<br/>SET status = 'failed',<br/>error_message = error,<br/>retry_count = retry_count + 1<br/>WHERE id = email_queue_id

            Consumer->>Kafka: Produce to email-retries topic
            Note over Kafka: { email_data, retry_count, next_retry_at: NOW() + 1hour }
        end
    end
```

## 5. Follow-up Scheduler

```mermaid
sequenceDiagram
    title Follow-up Scheduler

    participant Kafka as Kafka
    participant Scheduler as Follow-up Scheduler
    participant DB as Database

    Kafka->>Scheduler: Consume from followup-events topic
    Note over Kafka: { sequence_contact_id: id, step_completed: 0 }

    Scheduler->>DB: SELECT sc.*, s.steps<br/>FROM sequence_contacts sc<br/>JOIN sequences s ON sc.sequence_id = s.id<br/>WHERE sc.id = sequence_contact_id
    DB-->>Scheduler: Sequence contact with all steps

    Scheduler->>Scheduler: Calculate next step
    Note over Scheduler: current_step = step_completed + 1

    alt Has next step
        Scheduler->>DB: SELECT wait_days FROM steps<br/>WHERE sequence_id = sequence_id<br/>AND step_order = current_step
        DB-->>Scheduler: wait_days = 3

        Scheduler->>Scheduler: Calculate next send date
        Note over Scheduler: scheduled_for = NOW() + wait_days days

        Scheduler->>DB: BEGIN TRANSACTION
        Scheduler->>DB: UPDATE sequence_contacts<br/>SET current_step = current_step,<br/>next_send_at = scheduled_for<br/>WHERE id = sequence_contact_id

        Scheduler->>DB: INSERT INTO email_queues<br/>(sequence_contact_id, step_order,<br/>subject, content, scheduled_for, status)<br/>SELECT sequence_contact_id, current_step,<br/>s.subject, s.content, scheduled_for, 'scheduled'<br/>FROM steps s<br/>WHERE s.sequence_id = sequence_id<br/>AND s.step_order = current_step
        Scheduler->>DB: COMMIT TRANSACTION

    else No more steps
        Scheduler->>DB: UPDATE sequence_contacts<br/>SET status = 'completed',<br/>completed_at = NOW()<br/>WHERE id = sequence_contact_id
    end
```

## 6. Retry Mechanism

```mermaid
sequenceDiagram
    title Retry Mechanism

    participant Kafka as Kafka
    participant RetryConsumer as Retry Consumer
    participant DB as Database

    Note over RetryConsumer: Runs every hour

    Kafka->>RetryConsumer: Consume from email-retries topic
    Note over Kafka: { email_data, retry_count: 1, next_retry_at }

    RetryConsumer->>DB: Check if email still needs sending
    Note over DB: SELECT status, retry_count FROM email_queues WHERE id = email_queue_id
    DB-->>RetryConsumer: status = 'failed', retry_count = 1

    alt Under max retries (retry_count < 3)
        RetryConsumer->>DB: UPDATE email_queues<br/>SET status = 'scheduled',<br/>scheduled_for = NOW()<br/>WHERE id = email_queue_id

        RetryConsumer->>Kafka: Produce to email-jobs topic
        Note over Kafka: Single email for retry

    else Max retries exceeded (retry_count >= 3)
        RetryConsumer->>DB: UPDATE email_queues<br/>SET status = 'permanent_failure'<br/>WHERE id = email_queue_id

        RetryConsumer->>DB: UPDATE sequence_contacts<br/>SET status = 'bounced'<br/>WHERE id = sequence_contact_id
    end
```

## Database Schema (Accurate)

```sql
-- Core domain tables
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
```

## Kafka Topics Configuration

```yaml
topics:
  email-jobs:
    purpose: "Primary email processing queue"

  followup-events:
    purpose: "Trigger follow-up email scheduling"

  email-retries:
    purpose: "Failed email retries with exponential backoff"

  email-events:
    purpose: "Email tracking and analytics events"
```

## Highlights

- **Async Processing**: Non-blocking email sending with Kafka
- **Scalability**: Horizontal scaling of consumer instances
- **Fault Tolerance**: Automatic retries with exponential backoff
- **Capacity Management**: Respects mailbox daily limits
- **Monitoring**: Comprehensive tracking of all email events
- **Reliability**: Database transactions ensure data consistency
- **Performance**: Batch processing of emails per Kafka message
