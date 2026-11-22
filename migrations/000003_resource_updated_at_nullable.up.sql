ALTER TABLE graphics.resources
    ALTER COLUMN updated_at DROP NOT NULL,
    ALTER COLUMN updated_at DROP DEFAULT;