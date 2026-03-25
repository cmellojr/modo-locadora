-- =============================================================================
-- SEED: Dados iniciais para desenvolvimento
-- Jogos da Ação Games nº 1 (Julho 1991) + 3 sócios de teste
-- Executar via: go run ./cmd/server --seed
-- =============================================================================

-- Idempotência: só insere se não houver jogos
DO $seed$
BEGIN
    IF (SELECT COUNT(*) FROM games) > 0 THEN
        RAISE NOTICE 'Seed: banco ja populado. Pulando.';
        RETURN;
    END IF;

    -- ════════════════════════════════════════════════════════════════════════
    -- JOGOS — Ação Games nº 1 (Julho 1991)
    -- ════════════════════════════════════════════════════════════════════════

    -- Golden Axe (Mega Drive)
    INSERT INTO games (id, title, igdb_id, platform, summary, cover_url, source_magazine, acquired_at)
    VALUES ('a1b2c3d4-1111-4000-8000-000000000001', 'Golden Axe', '',
            'Mega Drive',
            'Hack and slash cooperativo no mundo de Yuria. Resgate a magia do machado dourado com Ax Battler, Tyris Flare ou Gilius Thunderhead!',
            'https://upload.wikimedia.org/wikipedia/en/a/ab/Golden_Axe_Promo_Flyer.jpeg',
            'Acao Games #1', '1991-07-01');
    INSERT INTO game_copies (id, game_id, status)
    VALUES ('c0010001-0001-4000-8000-000000000001', 'a1b2c3d4-1111-4000-8000-000000000001', 'available');

    -- Altered Beast (Mega Drive)
    INSERT INTO games (id, title, igdb_id, platform, summary, cover_url, source_magazine, acquired_at)
    VALUES ('a1b2c3d4-2222-4000-8000-000000000002', 'Altered Beast', '',
            'Mega Drive',
            'Rise from your grave! Luta mitologica lado a lado contra as forcas de Neff. O jogo que veio com o Mega Drive!',
            'https://upload.wikimedia.org/wikipedia/en/2/2b/Altered_Beast_Cover.jpg',
            'Acao Games #1', '1991-07-01');
    INSERT INTO game_copies (id, game_id, status)
    VALUES ('c0010001-0002-4000-8000-000000000002', 'a1b2c3d4-2222-4000-8000-000000000002', 'available');

    -- Super Mario Bros. 3 (NES / Nintendinho)
    INSERT INTO games (id, title, igdb_id, platform, summary, cover_url, source_magazine, acquired_at)
    VALUES ('a1b2c3d4-3333-4000-8000-000000000003', 'Super Mario Bros. 3', '',
            'NES',
            'O encanador mais famoso do mundo enfrenta os Koopalings em 8 mundos. Folha de Tanooki, sapo e martelo — o melhor Mario de todos!',
            'https://upload.wikimedia.org/wikipedia/en/a/a5/Super_Mario_Bros._3_coverart.png',
            'Acao Games #1', '1991-07-01');
    INSERT INTO game_copies (id, game_id, status)
    VALUES ('c0010001-0003-4000-8000-000000000003', 'a1b2c3d4-3333-4000-8000-000000000003', 'available');

    -- Castle of Illusion Starring Mickey Mouse (Mega Drive)
    INSERT INTO games (id, title, igdb_id, platform, summary, cover_url, source_magazine, acquired_at)
    VALUES ('a1b2c3d4-4444-4000-8000-000000000004', 'Castle of Illusion', '',
            'Mega Drive',
            'Mickey Mouse adentra o Castelo da Ilusao para salvar Minnie da bruxa Mizrabel. Plataforma magica da SEGA!',
            'https://upload.wikimedia.org/wikipedia/en/f/f4/Castle_of_illusion_Mickey_mouse.jpg',
            'Acao Games #1', '1991-07-01');
    INSERT INTO game_copies (id, game_id, status)
    VALUES ('c0010001-0004-4000-8000-000000000004', 'a1b2c3d4-4444-4000-8000-000000000004', 'available');

    -- Double Dragon II: The Revenge (NES / Nintendinho)
    INSERT INTO games (id, title, igdb_id, platform, summary, cover_url, source_magazine, acquired_at)
    VALUES ('a1b2c3d4-5555-4000-8000-000000000005', 'Double Dragon II: The Revenge', '',
            'NES',
            'Billy e Jimmy Lee vingam a morte de Marian neste classico beat em up cooperativo. Golpes devastadores e fases icônicas!',
            'https://upload.wikimedia.org/wikipedia/en/0/02/NES_Double_Dragon_II_packaging_front.jpg',
            'Acao Games #1', '1991-07-01');
    INSERT INTO game_copies (id, game_id, status)
    VALUES ('c0010001-0005-4000-8000-000000000005', 'a1b2c3d4-5555-4000-8000-000000000005', 'available');

    -- ════════════════════════════════════════════════════════════════════════
    -- SÓCIOS DE TESTE
    -- ════════════════════════════════════════════════════════════════════════
    -- Senhas: tio_da_locadora=sopre_a_fita | MegaDriveKid=sega1991 | Devedor=atrasado123 | Novato=novato2026

    -- Avança a sequence para 4 membros
    PERFORM nextval('membership_seq'); -- 1991-001
    PERFORM nextval('membership_seq'); -- 1991-002
    PERFORM nextval('membership_seq'); -- 1991-003
    PERFORM nextval('membership_seq'); -- 1991-004

    -- 0. Tio da Locadora — Administrador (email = ADMIN_EMAIL)
    INSERT INTO members (id, profile_name, email, password_hash, favorite_console, membership_number, joined_at)
    VALUES ('aabb0001-0000-4000-8000-000000000000',
            'tio_da_locadora', 'admin@locadora.com',
            '$2a$10$MNvVnff1TjnTrXDHY3OmiOKwx8NuVSgpMNaGXTEJ6soaJLxCSWJfG',
            'Mega Drive', '1991-001',
            '1991-07-01 08:00:00-03');

    -- 1. MegaDriveKid — Sócio exemplar, muitas devoluções no prazo
    INSERT INTO members (id, profile_name, email, password_hash, favorite_console, membership_number, joined_at)
    VALUES ('aabb0001-0001-4000-8000-000000000001',
            'MegaDriveKid', 'mega@locadora.com',
            '$2a$10$v/pOxtjrYzlrA5SbkO3EFubZN2tBWsZA4Fc673Fq8RMekVkSChyAO',
            'Mega Drive', '1991-002',
            '1991-07-15 10:00:00-03');

    -- 2. Devedor — Inadimplente, painel da vergonha
    INSERT INTO members (id, profile_name, email, password_hash, favorite_console, membership_number,
                         status, late_count, joined_at)
    VALUES ('aabb0001-0002-4000-8000-000000000002',
            'Devedor', 'devedor@locadora.com',
            '$2a$10$cYaEEdblvHr84QKT2c0toeZMIOUpgt4omo84FGaZAnJymY2s/inI.',
            'NES', '1991-003',
            'em_debito', 3,
            '1991-08-01 14:00:00-03');

    -- 3. Novato — Sócio novo, sem histórico
    INSERT INTO members (id, profile_name, email, password_hash, favorite_console, membership_number, joined_at)
    VALUES ('aabb0001-0003-4000-8000-000000000003',
            'Novato', 'novato@locadora.com',
            '$2a$10$mUiUVmj502aSoM5datTu9ukxCR/VS4IiEcAOFuC2eZkrnF.y.AWDa',
            'Mega Drive', '1991-004',
            '2026-03-10 09:00:00-03');

    -- ════════════════════════════════════════════════════════════════════════
    -- HISTÓRICO DE ALUGUÉIS
    -- ════════════════════════════════════════════════════════════════════════

    -- MegaDriveKid: 3 devoluções no prazo (exemplar!)
    -- Aluguel 1: Golden Axe — zerou
    INSERT INTO rentals (id, member_id, copy_id, rented_at, due_at, returned_at, public_legacy)
    VALUES ('ee000001-0001-4000-8000-000000000001',
            'aabb0001-0001-4000-8000-000000000001',
            'c0010001-0001-4000-8000-000000000001',
            NOW() - INTERVAL '30 days', NOW() - INTERVAL '27 days',
            NOW() - INTERVAL '28 days', 'zerei');

    -- Aluguel 2: Castle of Illusion — zerou
    INSERT INTO rentals (id, member_id, copy_id, rented_at, due_at, returned_at, public_legacy)
    VALUES ('ee000001-0002-4000-8000-000000000002',
            'aabb0001-0001-4000-8000-000000000001',
            'c0010001-0004-4000-8000-000000000004',
            NOW() - INTERVAL '20 days', NOW() - INTERVAL '17 days',
            NOW() - INTERVAL '18 days', 'zerei');

    -- Aluguel 3: Double Dragon II — jogou um pouco
    INSERT INTO rentals (id, member_id, copy_id, rented_at, due_at, returned_at, public_legacy)
    VALUES ('ee000001-0003-4000-8000-000000000003',
            'aabb0001-0001-4000-8000-000000000001',
            'c0010001-0005-4000-8000-000000000005',
            NOW() - INTERVAL '10 days', NOW() - INTERVAL '7 days',
            NOW() - INTERVAL '8 days', 'joguei_um_pouco');

    -- Devedor: 1 aluguel ativo e vencido (Altered Beast — há 10 dias, prazo há 7)
    INSERT INTO rentals (id, member_id, copy_id, rented_at, due_at)
    VALUES ('ee000001-0004-4000-8000-000000000004',
            'aabb0001-0002-4000-8000-000000000002',
            'c0010001-0002-4000-8000-000000000002',
            NOW() - INTERVAL '10 days', NOW() - INTERVAL '7 days');
    -- Marcar cópia como alugada
    UPDATE game_copies SET status = 'rented'
    WHERE id = 'c0010001-0002-4000-8000-000000000002';

    -- ════════════════════════════════════════════════════════════════════════
    -- FEED "ACONTECEU NA LOCADORA"
    -- ════════════════════════════════════════════════════════════════════════

    INSERT INTO activities (id, event_type, member_name, game_title, created_at)
    VALUES
        (gen_random_uuid(), 'new_game',          '', 'Golden Axe',                      NOW() - INTERVAL '5 days'),
        (gen_random_uuid(), 'new_game',          '', 'Altered Beast',                   NOW() - INTERVAL '5 days'),
        (gen_random_uuid(), 'new_game',          '', 'Super Mario Bros. 3',             NOW() - INTERVAL '5 days'),
        (gen_random_uuid(), 'new_game',          '', 'Castle of Illusion',              NOW() - INTERVAL '5 days'),
        (gen_random_uuid(), 'new_game',          '', 'Double Dragon II: The Revenge',   NOW() - INTERVAL '5 days'),
        (gen_random_uuid(), 'verdict_complete',  'MegaDriveKid', 'Golden Axe',          NOW() - INTERVAL '3 days'),
        (gen_random_uuid(), 'verdict_complete',  'MegaDriveKid', 'Castle of Illusion',  NOW() - INTERVAL '2 days'),
        (gen_random_uuid(), 'verdict_partial',   'MegaDriveKid', 'Double Dragon II: The Revenge', NOW() - INTERVAL '1 day'),
        (gen_random_uuid(), 'penalty',           'Devedor',      'Altered Beast',       NOW() - INTERVAL '6 hours');

    RAISE NOTICE 'Seed: concluido com sucesso!';
    RAISE NOTICE 'Socios: tio_da_locadora/sopre_a_fita (admin) | MegaDriveKid/sega1991 | Devedor/atrasado123 | Novato/novato2026';

END $seed$;
