-- 011_verdict_popularity.sql
-- Rename verdict slugs from Portuguese to English and update activity event types.

-- Verdict slugs: PT → EN
UPDATE rentals SET public_legacy = 'completed' WHERE public_legacy = 'zerei';
UPDATE rentals SET public_legacy = 'enjoyed' WHERE public_legacy = 'joguei_um_pouco';
UPDATE rentals SET public_legacy = 'gave_up' WHERE public_legacy = 'desisti';

-- Activity event types: old → new
UPDATE activities SET event_type = 'verdict_completed' WHERE event_type = 'verdict_complete';
UPDATE activities SET event_type = 'verdict_enjoyed' WHERE event_type = 'verdict_partial';
UPDATE activities SET event_type = 'verdict_gave_up' WHERE event_type = 'verdict_quit';
