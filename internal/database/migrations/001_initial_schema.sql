CREATE TABLE members (
    id UUID PRIMARY KEY,
    profile_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    favorite_console TEXT,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE games (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    igdb_id TEXT,
    platform TEXT NOT NULL,
    summary TEXT
);

CREATE TYPE game_copy_status AS ENUM ('available', 'rented');

CREATE TABLE game_copies (
    id UUID PRIMARY KEY,
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    status game_copy_status NOT NULL DEFAULT 'available'
);

CREATE TABLE rentals (
    id UUID PRIMARY KEY,
    member_id UUID NOT NULL REFERENCES members(id),
    copy_id UUID NOT NULL REFERENCES game_copies(id),
    rented_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    due_at TIMESTAMP WITH TIME ZONE NOT NULL,
    returned_at TIMESTAMP WITH TIME ZONE,
    personal_note TEXT,
    public_legacy TEXT
);
