CREATE TABLE IF NOT EXISTS redirect_events (
    timestamp DateTime,
    link text,
    alias text
)
ENGINE = MergeTree()
ORDER BY timestamp;