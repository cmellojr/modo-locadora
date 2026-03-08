-- Migration 005: Add member status and late_count for the auto-return reputation system.

-- status: 'active' (default) or 'em_debito' (in debt with the store).
ALTER TABLE members ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';

-- late_count: permanent counter of how many times a member returned late.
ALTER TABLE members ADD COLUMN IF NOT EXISTS late_count INTEGER NOT NULL DEFAULT 0;
