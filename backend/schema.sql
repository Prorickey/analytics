-- ALTER SYSTEM SET max_connections = 300;

CREATE TABLE IF NOT EXISTS analytics (
    id UUID DEFAULT gen_random_uuid(),
    event TEXT NOT NULL, 
    timestamp TIMESTAMP DEFAULT NOW(),
    metadata TEXT
);

CREATE TABLE IF NOT EXISTS analytics_minute (
    timestamp TIMESTAMP PRIMARY KEY,
    count BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS analytics_hour (
    timestamp TIMESTAMP PRIMARY KEY,
    count BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS analytics_day (
    timestamp TIMESTAMP PRIMARY KEY,
    count BIGINT DEFAULT 0
);

CREATE OR REPLACE FUNCTION increment_analytics_minute() 
RETURNS trigger as $$
BEGIN
    INSERT INTO analytics_minute(timestamp, count)
    VALUES(date_trunc('minute', NEW.timestamp), 1)
    ON CONFLICT (timestamp)
    DO UPDATE SET count = analytics_minute.count + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION increment_analytics_hour() 
RETURNS trigger as $$
BEGIN
    INSERT INTO analytics_hour(timestamp, count)
    VALUES(date_trunc('hour', NEW.timestamp), 1)
    ON CONFLICT (timestamp)
    DO UPDATE SET count = analytics_hour.count + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION increment_analytics_day() 
RETURNS trigger as $$
BEGIN
    INSERT INTO analytics_day(timestamp, count)
    VALUES(date_trunc('day', NEW.timestamp), 1)
    ON CONFLICT (timestamp)
    DO UPDATE SET count = analytics_day.count + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER analytics_insert_trigger_minute
AFTER INSERT ON analytics
FOR EACH ROW EXECUTE FUNCTION increment_analytics_minute();

CREATE OR REPLACE TRIGGER analytics_insert_trigger_hour
AFTER INSERT ON analytics
FOR EACH ROW EXECUTE FUNCTION increment_analytics_hour();

CREATE OR REPLACE TRIGGER analytics_insert_trigger_day
AFTER INSERT ON analytics
FOR EACH ROW EXECUTE FUNCTION increment_analytics_day();

CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    digest TEXT NOT NULL,
    token TEXT DEFAULT gen_random_uuid()
);

CREATE INDEX IF NOT EXISTS analytics_event_idx ON analytics(event);
CREATE INDEX IF NOT EXISTS analytics_timestamp_idx ON analytics(timestamp);
CREATE INDEX IF NOT EXISTS analytics_event_timestamp_idx ON analytics(event, timestamp);