CREATE TABLE IF NOT EXISTS events
(
    id           UUID DEFAULT generateUUIDv4(),
    short_url    String,
    original_url String,
    timestamp    DateTime,
    user_agent   String,
    ip           String
) ENGINE = MergeTree()
      ORDER BY id;
