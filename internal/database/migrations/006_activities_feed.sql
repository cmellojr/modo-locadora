-- Migration 006: Activities feed for "Aconteceu na Locadora".

CREATE TABLE IF NOT EXISTS activities (
    id          UUID PRIMARY KEY,
    event_type  TEXT NOT NULL,
    member_name TEXT NOT NULL DEFAULT '',
    game_title  TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_activities_created_at ON activities (created_at DESC);
