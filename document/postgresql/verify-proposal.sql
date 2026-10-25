-- schema-proposal.sql適用済みの空の検証DB向け。
-- psql -X -v ON_ERROR_STOP=1で実行する。fixtureは最後にROLLBACKする。
BEGIN;

INSERT INTO frontage.accounts(id, display_name) VALUES
    ('10000000-0000-0000-0000-000000000001', 'Alice'),
    ('10000000-0000-0000-0000-000000000002', 'Bob');
INSERT INTO frontage.account_identities(issuer, subject, account_id)
VALUES ('test-issuer', 'alice', '10000000-0000-0000-0000-000000000001');
INSERT INTO frontage.content_releases(id, version, engine_version, definition_schema_version, rules, manifest_sha256) VALUES
    ('20000000-0000-0000-0000-000000000001', 'test-r1', 'test-build', 1, '{}', decode(repeat('11', 32), 'hex')),
    ('20000000-0000-0000-0000-000000000002', 'test-r2', 'test-build', 1, '{}', decode(repeat('22', 32), 'hex'));
INSERT INTO frontage.card_definitions(id, card_key) VALUES
    ('30000000-0000-0000-0000-000000000001', 'unit.soldier'),
    ('30000000-0000-0000-0000-000000000002', 'spell.fireball');
INSERT INTO frontage.card_versions(release_id, card_id, kind, localization_key, definition) VALUES
    ('20000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'piece', 'cards.soldier', '{"hp":4}'),
    ('20000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000002', 'magic', 'cards.fireball', '{"damage":4}'),
    ('20000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000001', 'piece', 'cards.soldier', '{"hp":5}');

-- 版をまたぐカード種別ごとの所持枚数。初回配布と台帳を同じTXで記録する。
INSERT INTO frontage.account_cards(account_id, card_id) VALUES
    ('10000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001'),
    ('10000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000001');
INSERT INTO frontage.inventory_transactions(id, account_id, collection_revision, source_kind, source_id, request_sha256) VALUES
    ('80000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 1, 'starter_pack',
     '90000000-0000-0000-0000-000000000001', decode(repeat('55', 32), 'hex')),
    ('80000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000002', 1, 'starter_pack',
     '90000000-0000-0000-0000-000000000001', decode(repeat('55', 32), 'hex'));
INSERT INTO frontage.card_ownership_changes(transaction_id, account_id, card_id, delta, quantity_before, quantity_after) VALUES
    ('80000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 2, 0, 2),
    ('80000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000001', 2, 0, 2);
UPDATE frontage.account_cards SET quantity = 2;
UPDATE frontage.accounts SET collection_revision = 1;

-- headのFKは遅延されるので、ヘッダー→初版の順で作成できる。
INSERT INTO frontage.decks(id, account_id, name, current_revision) VALUES
    ('40000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 'Alice deck', 1),
    ('40000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000002', 'Bob deck', 1);
INSERT INTO frontage.deck_revisions(deck_id, revision, release_id) VALUES
    ('40000000-0000-0000-0000-000000000001', 1, '20000000-0000-0000-0000-000000000001'),
    ('40000000-0000-0000-0000-000000000002', 1, '20000000-0000-0000-0000-000000000001');
INSERT INTO frontage.deck_cards(deck_id, revision, release_id, zone, card_id, quantity) VALUES
    ('40000000-0000-0000-0000-000000000001', 1, '20000000-0000-0000-0000-000000000001', 'main', '30000000-0000-0000-0000-000000000001', 2),
    ('40000000-0000-0000-0000-000000000002', 1, '20000000-0000-0000-0000-000000000001', 'main', '30000000-0000-0000-0000-000000000001', 2);

INSERT INTO frontage.matches(id, release_id, mode, status, started_at, owner_node, fencing_token, lease_until) VALUES
    ('50000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', 'pve', 'active', now(), 'test-node', 1, now() + interval '1 minute'),
    ('50000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000001', 'pve', 'active', now(), 'test-node', 1, now() + interval '1 minute');
INSERT INTO frontage.match_participants(match_id, seat, release_id, engine_player_id, account_id, deck_id, deck_revision) VALUES
    ('50000000-0000-0000-0000-000000000001', 0, '20000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', '40000000-0000-0000-0000-000000000001', 1),
    ('50000000-0000-0000-0000-000000000002', 0, '20000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000001', '40000000-0000-0000-0000-000000000001', 1);
INSERT INTO frontage.match_participants(match_id, seat, release_id, engine_player_id, npc_profile) VALUES
    ('50000000-0000-0000-0000-000000000001', 1, '20000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000002', '{"strategy":"test","version":1,"deck":[]}'),
    ('50000000-0000-0000-0000-000000000002', 1, '20000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000004', '{"strategy":"test","version":1,"deck":[]}');
-- stateの意味的な完全性はアプリ側の検査対象。ここではSQL制約だけを検証する。
INSERT INTO frontage.match_snapshots(match_id, seq, state_schema_version, state) VALUES
    ('50000000-0000-0000-0000-000000000001', 0, 1, '{}'),
    ('50000000-0000-0000-0000-000000000002', 0, 1, '{}');
INSERT INTO frontage.active_match_slots(account_id, match_id, seat)
VALUES ('10000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 0);

INSERT INTO frontage.match_snapshots(match_id, seq, state_schema_version, state)
VALUES ('50000000-0000-0000-0000-000000000001', 1, 1, '{"test_seq":1}');
INSERT INTO frontage.match_commands(match_id, seq, command_id, actor_seat, expected_seq, protocol_version,
                                    command_kind, request_sha256, payload, response)
VALUES ('50000000-0000-0000-0000-000000000001', 1, '70000000-0000-0000-0000-000000000001', 0, 0, 1,
        'end_turn', decode(repeat('33', 32), 'hex'), '{}', '{"accepted_seq":1}');
UPDATE frontage.matches SET last_seq = 1 WHERE id = '50000000-0000-0000-0000-000000000001';

-- ROLLBACK任せにせず、fixtureの遅延FKもこの時点で検証する。
SET CONSTRAINTS ALL IMMEDIATE;

DO $$
BEGIN
    BEGIN
        UPDATE frontage.deck_cards SET quantity = 0;
        RAISE EXCEPTION 'zero quantity was accepted';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.deck_cards SET zone = 'hand';
        RAISE EXCEPTION 'hand was accepted as a deck-building zone';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.deck_cards SET release_id = '20000000-0000-0000-0000-000000000002';
        RAISE EXCEPTION 'mixed release was accepted';
    EXCEPTION WHEN foreign_key_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.match_participants SET deck_id = '40000000-0000-0000-0000-000000000002'
        WHERE match_id = '50000000-0000-0000-0000-000000000001' AND seat = 0;
        RAISE EXCEPTION 'another account deck was accepted';
    EXCEPTION WHEN foreign_key_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.matches SET last_seq = 99 WHERE id = '50000000-0000-0000-0000-000000000001';
        RAISE EXCEPTION 'head without snapshot was accepted';
    EXCEPTION WHEN foreign_key_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.decks SET current_revision = 99;
        RAISE EXCEPTION 'deck head without revision was accepted';
    EXCEPTION WHEN foreign_key_violation THEN NULL;
    END;
    BEGIN
        INSERT INTO frontage.active_match_slots(account_id, match_id, seat)
        VALUES ('10000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000002', 0);
        RAISE EXCEPTION 'two active matches for one account were accepted';
    EXCEPTION WHEN unique_violation THEN NULL;
    END;
    BEGIN
        INSERT INTO frontage.match_snapshots(match_id, seq, state_schema_version, state)
        VALUES ('50000000-0000-0000-0000-000000000001', 2, 1, '{}');
        INSERT INTO frontage.match_commands(match_id, seq, command_id, actor_seat, expected_seq, protocol_version,
                                            command_kind, request_sha256, payload, response)
        VALUES ('50000000-0000-0000-0000-000000000001', 2, '70000000-0000-0000-0000-000000000001', 0, 1, 1,
                'end_turn', decode(repeat('33', 32), 'hex'), '{}', '{}');
        RAISE EXCEPTION 'duplicate command id was accepted';
    EXCEPTION WHEN unique_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.match_commands SET expected_seq = 9;
        RAISE EXCEPTION 'invalid expected sequence was accepted';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.match_snapshots SET state = '[]';
        RAISE EXCEPTION 'non-object snapshot was accepted';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.match_participants SET npc_profile = '{}' WHERE account_id IS NOT NULL;
        RAISE EXCEPTION 'participant was both human and NPC';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.matches SET status = 'finished', ended_at = now(), outcome = 'win',
            winner_seat = NULL, end_reason = 'test', owner_node = NULL, lease_until = NULL
        WHERE id = '50000000-0000-0000-0000-000000000001';
        RAISE EXCEPTION 'win without winner was accepted';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.matches SET status = 'finished', ended_at = now(), outcome = NULL,
            end_reason = 'test', owner_node = NULL, lease_until = NULL
        WHERE id = '50000000-0000-0000-0000-000000000001';
        RAISE EXCEPTION 'finished match without outcome was accepted';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.matches SET history_status = 'archived', replay_object_key = 'test/replay';
        RAISE EXCEPTION 'active match was archived';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        INSERT INTO frontage.account_identities(issuer, subject, account_id)
        VALUES ('test-issuer', 'alice', '10000000-0000-0000-0000-000000000002');
        RAISE EXCEPTION 'one external identity was linked to two accounts';
    EXCEPTION WHEN unique_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.match_participants SET release_id = '20000000-0000-0000-0000-000000000002'
        WHERE match_id = '50000000-0000-0000-0000-000000000001' AND seat = 0;
        RAISE EXCEPTION 'participant release did not match its game';
    EXCEPTION WHEN foreign_key_violation THEN NULL;
    END;
    BEGIN
        INSERT INTO frontage.match_commands(match_id, seq, command_id, actor_seat, expected_seq, protocol_version,
                                            command_kind, request_sha256, payload, response)
        VALUES ('50000000-0000-0000-0000-000000000001', 99, '70000000-0000-0000-0000-000000000099', 0, 98, 1,
                'end_turn', decode(repeat('33', 32), 'hex'), '{}', '{}');
        RAISE EXCEPTION 'command without snapshot was accepted';
    EXCEPTION WHEN foreign_key_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.matches SET fencing_token = 0
        WHERE id = '50000000-0000-0000-0000-000000000001';
        RAISE EXCEPTION 'owner was assigned without a positive fencing token';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.account_cards SET quantity = -1;
        RAISE EXCEPTION 'negative ownership balance was accepted';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        INSERT INTO frontage.inventory_transactions(id, account_id, collection_revision, source_kind, source_id, request_sha256)
        VALUES ('80000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000001', 2,
                'starter_pack', '90000000-0000-0000-0000-000000000001', decode(repeat('55', 32), 'hex'));
        RAISE EXCEPTION 'one ownership source was applied twice';
    EXCEPTION WHEN unique_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.card_ownership_changes SET delta = 0;
        RAISE EXCEPTION 'zero ownership change was accepted';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.card_ownership_changes SET account_id = '10000000-0000-0000-0000-000000000002'
        WHERE transaction_id = '80000000-0000-0000-0000-000000000001';
        RAISE EXCEPTION 'ownership change used another account transaction';
    EXCEPTION WHEN foreign_key_violation THEN NULL;
    END;
    BEGIN
        UPDATE frontage.card_ownership_changes SET delta = 3;
        RAISE EXCEPTION 'ownership delta disagreed with before/after quantities';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
    RAISE NOTICE '23 invalid writes rejected';
END $$;

-- 正常な終了・参加枠解放も同じトランザクションで成立する。
INSERT INTO frontage.match_snapshots(match_id, seq, state_schema_version, state)
VALUES ('50000000-0000-0000-0000-000000000001', 2, 1, '{"test_seq":2,"over":true}');
INSERT INTO frontage.match_commands(match_id, seq, command_id, expected_seq, protocol_version,
                                    command_kind, request_sha256, payload, response)
VALUES ('50000000-0000-0000-0000-000000000001', 2, '70000000-0000-0000-0000-000000000002', 1, 1,
        'finish_match', decode(repeat('44', 32), 'hex'), '{}', '{"accepted_seq":2}');
UPDATE frontage.matches SET last_seq = 2, status = 'finished', ended_at = now(), outcome = 'win', winner_seat = 0,
    end_reason = 'test', owner_node = NULL, lease_until = NULL
WHERE id = '50000000-0000-0000-0000-000000000001';
DELETE FROM frontage.active_match_slots WHERE match_id = '50000000-0000-0000-0000-000000000001';
SET CONSTRAINTS ALL IMMEDIATE;

-- アーカイブ後は開始・最終状態を残し、中間の操作と状態を同じTXで削除できる。
-- オブジェクト保存自体は実装対象外。このキーは制約検証用fixture。
UPDATE frontage.matches SET history_status = 'archived', replay_object_key = 'test/complete-replay'
WHERE id = '50000000-0000-0000-0000-000000000001';
SET CONSTRAINTS ALL DEFERRED;
DELETE FROM frontage.match_snapshots WHERE match_id = '50000000-0000-0000-0000-000000000001' AND seq = 1;
DELETE FROM frontage.match_commands WHERE match_id = '50000000-0000-0000-0000-000000000001' AND seq = 1;
SET CONSTRAINTS ALL IMMEDIATE;
SELECT 'schema proposal checks passed' AS result;
ROLLBACK;
