CREATE TABLE IF NOT EXISTS matches (
    id         TEXT PRIMARY KEY,
    data       JSONB        NOT NULL,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- Optional: index for waiting lobby lookups by status inside JSON
CREATE INDEX IF NOT EXISTS matches_waiting_idx
    ON matches ((data ->> 'Status'))
    WHERE (data ->> 'Status') = 'in_progress' OR (data ->> 'Status') = 'waiting';


