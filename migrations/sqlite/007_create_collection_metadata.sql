-- Create collection_metadata table to track incremental collection state
CREATE TABLE IF NOT EXISTS collection_metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repository TEXT UNIQUE NOT NULL,
    last_collected_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast lookups by repository
CREATE INDEX IF NOT EXISTS idx_collection_metadata_repository
    ON collection_metadata(repository);

-- Trigger to auto-update updated_at
CREATE TRIGGER IF NOT EXISTS trg_collection_metadata_updated_at
BEFORE UPDATE ON collection_metadata
FOR EACH ROW
BEGIN
    UPDATE collection_metadata SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
