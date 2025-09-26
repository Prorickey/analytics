CREATE TABLE IF NOT EXISTS analytics (
    id UUID DEFAULT gen_random_uuid(),
    event TEXT NOT NULL, 
    timestamp TIMESTAMP DEFAULT NOW(),
    metadata TEXT
);

CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    digest TEXT NOT NULL,
    token TEXT DEFAULT gen_random_uuid()
);

CREATE INDEX IF NOT EXISTS analytics_event_idx ON analytics(event);
CREATE INDEX IF NOT EXISTS analytics_timestamp_idx ON analytics(timestamp);
CREATE INDEX IF NOT EXISTS analytics_event_timestamp_idx ON analytics(event, timestamp);