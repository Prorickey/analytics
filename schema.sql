-- ALTER SYSTEM SET max_connections = 300;

CREATE TABLE IF NOT EXISTS analytics (
    id UUID DEFAULT gen_random_uuid(),
    event TEXT NOT NULL, 
    timestamp TIMESTAMP DEFAULT NOW(),
    metadata TEXT
);

CREATE TABLE IF NOT EXISTS analytics_minute (
    event TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    count BIGINT DEFAULT 0,

    PRIMARY KEY (event, timestamp)
);

CREATE TABLE IF NOT EXISTS analytics_hour (
    event TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    count BIGINT DEFAULT 0,

    PRIMARY KEY (event, timestamp)
);

CREATE TABLE IF NOT EXISTS analytics_day (
    event TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    count BIGINT DEFAULT 0,

    PRIMARY KEY (event, timestamp)
);

CREATE OR REPLACE FUNCTION increment_analytics_minute() 
RETURNS trigger as $$
BEGIN
    INSERT INTO analytics_minute(event, timestamp, count)
    VALUES(NEW.event, date_trunc('minute', NEW.timestamp), 1)
    ON CONFLICT (event, timestamp)
    DO UPDATE SET count = analytics_minute.count + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION increment_analytics_hour() 
RETURNS trigger as $$
BEGIN
    INSERT INTO analytics_hour(event, timestamp, count)
    VALUES(NEW.event, date_trunc('hour', NEW.timestamp), 1)
    ON CONFLICT (event, timestamp)
    DO UPDATE SET count = analytics_hour.count + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION increment_analytics_day() 
RETURNS trigger as $$
BEGIN
    INSERT INTO analytics_day(event, timestamp, count)
    VALUES(NEW.event, date_trunc('day', NEW.timestamp), 1)
    ON CONFLICT (event, timestamp)
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

CREATE TABLE IF NOT EXISTS service_auth (
    id UUID DEFAULT gen_random_uuid(),
    service TEXT NOT NULL,
    token TEXT DEFAULT gen_random_uuid()
);

CREATE INDEX IF NOT EXISTS analytics_event_idx ON analytics(event);
CREATE INDEX IF NOT EXISTS analytics_timestamp_idx ON analytics(timestamp);
CREATE INDEX IF NOT EXISTS analytics_event_timestamp_idx ON analytics(event, timestamp);

CREATE INDEX IF NOT EXISTS analytics_minute_event_idx ON analytics_minute(event);
CREATE INDEX IF NOT EXISTS analytics_minute_timestamp_idx ON analytics_minute(timestamp);
CREATE INDEX IF NOT EXISTS analytics_minute_event_timestamp_idx ON analytics_minute(event, timestamp);

CREATE INDEX IF NOT EXISTS analytics_hour_event_idx ON analytics_hour(event);
CREATE INDEX IF NOT EXISTS analytics_hour_timestamp_idx ON analytics_hour(timestamp);
CREATE INDEX IF NOT EXISTS analytics_hour_event_timestamp_idx ON analytics_hour(event, timestamp);

CREATE INDEX IF NOT EXISTS analytics_day_event_idx ON analytics_day(event);
CREATE INDEX IF NOT EXISTS analytics_day_timestamp_idx ON analytics_day(timestamp);
CREATE INDEX IF NOT EXISTS analytics_day_event_timestamp_idx ON analytics_day(event, timestamp);