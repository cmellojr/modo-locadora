-- Migration 003: Add membership fields, ensure game copies exist, seed admin member.

-- 1. Add new columns to members table.
ALTER TABLE members ADD COLUMN IF NOT EXISTS membership_number TEXT UNIQUE;
ALTER TABLE members ADD COLUMN IF NOT EXISTS address TEXT DEFAULT '';
ALTER TABLE members ADD COLUMN IF NOT EXISTS phone TEXT DEFAULT '';

-- 2. Create a sequence for membership numbers (starting at 1).
CREATE SEQUENCE IF NOT EXISTS membership_seq START 1;

-- 3. Backfill membership_number for existing members that don't have one.
UPDATE members
SET membership_number = '1991-' || LPAD(nextval('membership_seq')::TEXT, 3, '0')
WHERE membership_number IS NULL;

-- 4. Create game_copies for any games that don't have a copy yet.
--    Each game gets one physical copy (cartridge) by default.
INSERT INTO game_copies (id, game_id, status)
SELECT gen_random_uuid(), g.id, 'available'
FROM games g
WHERE NOT EXISTS (
    SELECT 1 FROM game_copies gc WHERE gc.game_id = g.id
);
