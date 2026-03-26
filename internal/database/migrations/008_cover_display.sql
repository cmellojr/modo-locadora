-- 008: Add cover_display column for controlling how cover images are rendered.
-- Values: 'cover' (default, crops to fill), 'contain' (show full image), 'fill' (stretch).
ALTER TABLE games ADD COLUMN IF NOT EXISTS cover_display VARCHAR(10) NOT NULL DEFAULT 'cover';
