-- PostgreSQL 18向けの設計確認用DDL。アプリのマイグレーションではない。
-- 空の検証DBにのみ適用する。UUIDはアプリケーションで生成する。
BEGIN;
CREATE SCHEMA frontage;
REVOKE ALL ON SCHEMA frontage FROM PUBLIC;

CREATE TABLE frontage.accounts (
    id uuid PRIMARY KEY,
    display_name text NOT NULL CHECK (char_length(btrim(display_name)) BETWEEN 1 AND 80),
    status text NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'deleted')),
    collection_revision bigint NOT NULL DEFAULT 0 CHECK (collection_revision >= 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE frontage.account_identities (
    issuer text NOT NULL CHECK (issuer <> ''),
    subject text NOT NULL CHECK (subject <> ''),
    account_id uuid NOT NULL REFERENCES frontage.accounts(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (issuer, subject)
);
CREATE INDEX account_identities_account_idx ON frontage.account_identities(account_id);

-- 公開済みのカード・ルール・エンジンの組を固定する。
-- ランタイムのDBロールは本テーブルとカード定義をSELECTのみ許可する。
CREATE TABLE frontage.content_releases (
    id uuid PRIMARY KEY,
    version text NOT NULL UNIQUE CHECK (version <> ''),
    engine_version text NOT NULL CHECK (engine_version <> ''),
    definition_schema_version integer NOT NULL CHECK (definition_schema_version > 0),
    rules jsonb NOT NULL CHECK (jsonb_typeof(rules) = 'object'),
    manifest_sha256 bytea NOT NULL CHECK (octet_length(manifest_sha256) = 32),
    published_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE frontage.card_definitions (
    id uuid PRIMARY KEY,
    card_key text NOT NULL UNIQUE CHECK (card_key ~ '^[a-z][a-z0-9_.-]*$'),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE frontage.card_versions (
    release_id uuid NOT NULL REFERENCES frontage.content_releases(id) ON DELETE RESTRICT,
    card_id uuid NOT NULL REFERENCES frontage.card_definitions(id) ON DELETE RESTRICT,
    kind text NOT NULL CHECK (kind IN ('piece', 'charm', 'structure', 'magic', 'skill')),
    localization_key text NOT NULL CHECK (localization_key <> ''),
    definition jsonb NOT NULL CHECK (jsonb_typeof(definition) = 'object'),
    PRIMARY KEY (release_id, card_id)
);
CREATE INDEX card_versions_card_idx ON frontage.card_versions(card_id, release_id);

-- 所持はカード種別ごとの枚数。release変更や対戦内個体IDから独立させる。
CREATE TABLE frontage.account_cards (
    account_id uuid NOT NULL REFERENCES frontage.accounts(id) ON DELETE RESTRICT,
    card_id uuid NOT NULL REFERENCES frontage.card_definitions(id) ON DELETE RESTRICT,
    quantity integer NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (account_id, card_id)
);

-- 1回の付与・消費で複数種類のカードを扱える。要求全体のハッシュを保存する。
CREATE TABLE frontage.inventory_transactions (
    id uuid PRIMARY KEY,
    account_id uuid NOT NULL REFERENCES frontage.accounts(id) ON DELETE RESTRICT,
    collection_revision bigint NOT NULL CHECK (collection_revision > 0),
    source_kind text NOT NULL CHECK (source_kind <> ''),
    source_id uuid NOT NULL,
    request_sha256 bytea NOT NULL CHECK (octet_length(request_sha256) = 32),
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (account_id, source_kind, source_id),
    UNIQUE (account_id, collection_revision),
    UNIQUE (id, account_id)
);

CREATE TABLE frontage.card_ownership_changes (
    transaction_id uuid NOT NULL,
    account_id uuid NOT NULL,
    card_id uuid NOT NULL,
    delta integer NOT NULL CHECK (delta <> 0),
    quantity_before integer NOT NULL CHECK (quantity_before >= 0),
    quantity_after integer NOT NULL CHECK (quantity_after >= 0),
    PRIMARY KEY (transaction_id, card_id),
    CHECK (quantity_after - quantity_before = delta),
    FOREIGN KEY (transaction_id, account_id)
        REFERENCES frontage.inventory_transactions(id, account_id) ON DELETE RESTRICT,
    FOREIGN KEY (account_id, card_id)
        REFERENCES frontage.account_cards(account_id, card_id) ON DELETE RESTRICT
);
CREATE INDEX card_ownership_changes_card_idx
    ON frontage.card_ownership_changes(account_id, card_id, transaction_id);

CREATE TABLE frontage.decks (
    id uuid PRIMARY KEY,
    account_id uuid NOT NULL REFERENCES frontage.accounts(id) ON DELETE RESTRICT,
    name text NOT NULL CHECK (char_length(btrim(name)) BETWEEN 1 AND 100),
    current_revision integer NOT NULL CHECK (current_revision > 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,
    UNIQUE (id, account_id)
);
CREATE INDEX decks_owner_recent_idx ON frontage.decks(account_id, updated_at DESC, id)
    WHERE deleted_at IS NULL;

CREATE TABLE frontage.deck_revisions (
    deck_id uuid NOT NULL REFERENCES frontage.decks(id) ON DELETE CASCADE,
    revision integer NOT NULL CHECK (revision > 0),
    release_id uuid NOT NULL REFERENCES frontage.content_releases(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (deck_id, revision),
    UNIQUE (deck_id, revision, release_id)
);
ALTER TABLE frontage.decks ADD CONSTRAINT decks_current_revision_fk
    FOREIGN KEY (id, current_revision) REFERENCES frontage.deck_revisions(deck_id, revision)
    DEFERRABLE INITIALLY DEFERRED;

-- hand/graveyardは構築情報ではなく、対戦スナップショットに保存する。
CREATE TABLE frontage.deck_cards (
    deck_id uuid NOT NULL,
    revision integer NOT NULL,
    release_id uuid NOT NULL,
    zone text NOT NULL CHECK (zone IN ('main', 'sub', 'extra')),
    card_id uuid NOT NULL,
    quantity integer NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (deck_id, revision, zone, card_id),
    FOREIGN KEY (deck_id, revision, release_id)
        REFERENCES frontage.deck_revisions(deck_id, revision, release_id) ON DELETE CASCADE,
    FOREIGN KEY (release_id, card_id)
        REFERENCES frontage.card_versions(release_id, card_id) ON DELETE RESTRICT
);

CREATE TABLE frontage.matches (
    id uuid PRIMARY KEY,
    release_id uuid NOT NULL REFERENCES frontage.content_releases(id) ON DELETE RESTRICT,
    mode text NOT NULL CHECK (mode IN ('pve', 'pvp')),
    status text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'finished', 'aborted')),
    last_seq bigint NOT NULL DEFAULT 0 CHECK (last_seq >= 0),
    outcome text CHECK (outcome IN ('win', 'draw', 'aborted')),
    winner_seat smallint CHECK (winner_seat BETWEEN 0 AND 1),
    end_reason text,
    owner_node text CHECK (owner_node IS NULL OR owner_node <> ''),
    fencing_token bigint NOT NULL DEFAULT 0 CHECK (fencing_token >= 0),
    lease_until timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    started_at timestamptz,
    ended_at timestamptz,
    history_status text NOT NULL DEFAULT 'complete' CHECK (history_status IN ('complete', 'archived', 'pruned')),
    replay_object_key text,
    UNIQUE (id, release_id),
    CHECK (started_at IS NULL OR started_at >= created_at),
    CHECK (ended_at IS NULL OR ended_at >= created_at),
    CHECK (started_at IS NULL OR ended_at IS NULL OR ended_at >= started_at),
    CHECK (status <> 'pending' OR started_at IS NULL),
    CHECK (status NOT IN ('active', 'finished') OR started_at IS NOT NULL),
    CHECK ((status IN ('pending', 'active')) = (ended_at IS NULL)),
    CHECK (
        (ended_at IS NULL AND outcome IS NULL AND winner_seat IS NULL AND end_reason IS NULL)
        OR
        (ended_at IS NOT NULL AND outcome IS NOT NULL AND end_reason IS NOT NULL AND end_reason <> ''
         AND ((status = 'finished' AND ((outcome = 'win' AND winner_seat IS NOT NULL)
                                       OR (outcome = 'draw' AND winner_seat IS NULL)))
              OR (status = 'aborted' AND outcome = 'aborted' AND winner_seat IS NULL)))
    ),
    CHECK ((owner_node IS NULL) = (lease_until IS NULL)),
    CHECK (owner_node IS NULL OR (status = 'active' AND fencing_token > 0)),
    CHECK (history_status <> 'archived' OR (replay_object_key IS NOT NULL AND replay_object_key <> '')),
    CHECK (status NOT IN ('pending', 'active') OR history_status = 'complete')
);
CREATE INDEX matches_expired_lease_idx ON frontage.matches(lease_until, id)
    WHERE status = 'active';
CREATE INDEX matches_ended_idx ON frontage.matches(ended_at, id)
    WHERE status IN ('finished', 'aborted');

CREATE TABLE frontage.match_participants (
    match_id uuid NOT NULL,
    seat smallint NOT NULL CHECK (seat BETWEEN 0 AND 1),
    release_id uuid NOT NULL,
    engine_player_id uuid NOT NULL UNIQUE,
    account_id uuid REFERENCES frontage.accounts(id) ON DELETE RESTRICT,
    deck_id uuid,
    deck_revision integer,
    npc_profile jsonb,
    PRIMARY KEY (match_id, seat),
    UNIQUE (match_id, account_id),
    UNIQUE (match_id, seat, account_id),
    FOREIGN KEY (match_id, release_id)
        REFERENCES frontage.matches(id, release_id) ON DELETE CASCADE,
    FOREIGN KEY (deck_id, account_id)
        REFERENCES frontage.decks(id, account_id) ON DELETE RESTRICT,
    FOREIGN KEY (deck_id, deck_revision, release_id)
        REFERENCES frontage.deck_revisions(deck_id, revision, release_id) ON DELETE RESTRICT,
    CHECK (
        (account_id IS NOT NULL AND deck_id IS NOT NULL AND deck_revision IS NOT NULL AND npc_profile IS NULL)
        OR
        (account_id IS NULL AND deck_id IS NULL AND deck_revision IS NULL
         AND npc_profile IS NOT NULL AND jsonb_typeof(npc_profile) = 'object')
    )
);
CREATE INDEX match_participants_account_idx ON frontage.match_participants(account_id, match_id)
    WHERE account_id IS NOT NULL;
ALTER TABLE frontage.matches ADD CONSTRAINT matches_winner_fk
    FOREIGN KEY (id, winner_seat) REFERENCES frontage.match_participants(match_id, seat)
    DEFERRABLE INITIALLY DEFERRED;

CREATE TABLE frontage.match_snapshots (
    match_id uuid NOT NULL REFERENCES frontage.matches(id) ON DELETE CASCADE,
    seq bigint NOT NULL CHECK (seq >= 0),
    state_schema_version integer NOT NULL CHECK (state_schema_version > 0),
    state jsonb NOT NULL CHECK (jsonb_typeof(state) = 'object'),
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (match_id, seq)
);
ALTER TABLE frontage.matches ADD CONSTRAINT matches_head_snapshot_fk
    FOREIGN KEY (id, last_seq) REFERENCES frontage.match_snapshots(match_id, seq)
    DEFERRABLE INITIALLY DEFERRED;

CREATE TABLE frontage.match_commands (
    match_id uuid NOT NULL REFERENCES frontage.matches(id) ON DELETE CASCADE,
    seq bigint NOT NULL CHECK (seq > 0),
    command_id uuid NOT NULL,
    actor_seat smallint CHECK (actor_seat BETWEEN 0 AND 1),
    expected_seq bigint NOT NULL CHECK (expected_seq >= 0),
    protocol_version integer NOT NULL CHECK (protocol_version > 0),
    command_kind text NOT NULL CHECK (command_kind <> ''),
    request_sha256 bytea NOT NULL CHECK (octet_length(request_sha256) = 32),
    payload jsonb NOT NULL CHECK (jsonb_typeof(payload) = 'object'),
    response jsonb NOT NULL CHECK (jsonb_typeof(response) = 'object'),
    committed_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (match_id, seq),
    UNIQUE (match_id, command_id),
    CHECK (expected_seq = seq - 1),
    FOREIGN KEY (match_id, actor_seat)
        REFERENCES frontage.match_participants(match_id, seat),
    FOREIGN KEY (match_id, seq)
        REFERENCES frontage.match_snapshots(match_id, seq)
        DEFERRABLE INITIALLY DEFERRED
);

-- 1アカウント1進行中対戦、を初期案とする。終了時に同じTXで削除する。
CREATE TABLE frontage.active_match_slots (
    account_id uuid PRIMARY KEY,
    match_id uuid NOT NULL,
    seat smallint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (match_id, seat, account_id)
        REFERENCES frontage.match_participants(match_id, seat, account_id) ON DELETE CASCADE
);
CREATE INDEX active_match_slots_match_idx ON frontage.active_match_slots(match_id);

COMMIT;
