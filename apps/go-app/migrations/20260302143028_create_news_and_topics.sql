-- +goose Up
-- +goose StatementBegin
CREATE TYPE news_status AS ENUM ('DRAFT', 'PUBLISHED', 'DELETED');

CREATE TABLE news (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL,
    content TEXT NOT NULL,
    status news_status NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_news_status ON news(status);
CREATE INDEX idx_news_deleted_at ON news(deleted_at);

CREATE TABLE topic (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_topic_deleted_at ON topic(deleted_at);

CREATE TABLE news_topics (
    news_id UUID NOT NULL REFERENCES news(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES topic(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (news_id, topic_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS news_topics;
DROP TABLE IF EXISTS topic;
DROP TABLE IF EXISTS news;
DROP TYPE IF EXISTS news_status;
-- +goose StatementEnd
