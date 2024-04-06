-- Create videos table
CREATE TABLE IF NOT EXISTS videos (
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    published_at TIMESTAMP,
    thumbnails TEXT
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_published_at ON videos(published_at);
CREATE INDEX IF NOT EXISTS idx_title ON videos(title);
CREATE INDEX IF NOT EXISTS idx_description ON videos(description);
