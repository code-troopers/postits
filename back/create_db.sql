CREATE TABLE users (
    id UUID PRIMARY KEY,
	given_name  TEXT NOT NULL,
	family_name TEXT NOT NULL,
	email      TEXT NOT NULL,
	picture    TEXT NOT NULL
);

CREATE TABLE boards (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE postits (
    id UUID PRIMARY KEY,
    board_id UUID REFERENCES boards(id),
    text TEXT,
    pos_x INT,
    pos_y INT,
    votes INT DEFAULT 0,
    show BOOLEAN DEFAULT TRUE,
    author_id UUID REFERENCES users(id),
    weight INT DEFAULT 0
);
