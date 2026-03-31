-- Migration 009: Clubs (turmas) support.
-- Adds clubs and club_members tables for the first M2M relationship.

CREATE TABLE IF NOT EXISTS clubs (
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    badge_url   TEXT NOT NULL DEFAULT '',
    website_url TEXT NOT NULL DEFAULT '',
    created_by  UUID NOT NULL REFERENCES members(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS club_members (
    club_id   UUID NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    member_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    role      TEXT NOT NULL DEFAULT 'member',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (club_id, member_id)
);

CREATE INDEX IF NOT EXISTS idx_club_members_member ON club_members(member_id);
CREATE INDEX IF NOT EXISTS idx_clubs_created_at ON clubs(created_at DESC);

-- Seed club data (only if no clubs exist yet).
DO $club_seed$
BEGIN
    IF (SELECT COUNT(*) FROM clubs) > 0 THEN
        RETURN;
    END IF;

    INSERT INTO clubs (id, name, description, badge_url, website_url, created_by, created_at)
    VALUES ('bb000001-0001-4000-8000-000000000001',
            'Turma da Acao Games',
            'Galera que cresceu lendo a revista Acao Games e trocando fitas na locadora.',
            '', '',
            'aabb0001-0001-4000-8000-000000000001',
            NOW());

    INSERT INTO club_members (club_id, member_id, role, joined_at) VALUES
        ('bb000001-0001-4000-8000-000000000001', 'aabb0001-0001-4000-8000-000000000001', 'admin', NOW()),
        ('bb000001-0001-4000-8000-000000000001', 'aabb0001-0003-4000-8000-000000000003', 'member', NOW());

END $club_seed$;
