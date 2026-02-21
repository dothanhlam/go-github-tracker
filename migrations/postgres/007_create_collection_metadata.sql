-- Create collection_metadata table to track incremental collection state
CREATE TABLE IF NOT EXISTS collection_metadata (
    id SERIAL PRIMARY KEY,
    repository VARCHAR(255) UNIQUE NOT NULL,
    last_collected_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Index for fast lookups by repository
CREATE INDEX IF NOT EXISTS idx_collection_metadata_repository
    ON collection_metadata(repository);

-- Trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_collection_metadata_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_collection_metadata_updated_at
BEFORE UPDATE ON collection_metadata
FOR EACH ROW
EXECUTE FUNCTION update_collection_metadata_updated_at();
