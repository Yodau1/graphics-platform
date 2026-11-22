UPDATE graphics.resources
SET updated_at = created_at
WHERE updated_at IS NULL;

ALTER TABLE graphics.resources
    ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP,
    ALTER COLUMN updated_at SET NOT NULL;