-- Migration 010: Rename member status 'em_debito' to 'in_debt' (English convention).
UPDATE members SET status = 'in_debt' WHERE status = 'em_debito';
