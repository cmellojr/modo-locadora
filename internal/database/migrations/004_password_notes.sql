-- Migration 004: Add password_notes column for the "Caderno de Passwords" feature.
ALTER TABLE members ADD COLUMN IF NOT EXISTS password_notes TEXT DEFAULT '';
