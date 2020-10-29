CREATE TABLE IF NOT EXISTS profile
(
    id       serial PRIMARY KEY,
    nickname text UNIQUE,
    -- Данное поле допускает только латиницу, цифры и знак подчеркивания.
    -- Сравнение имени регистронезависимо
    fullname text        NOT NULL,
    about    text,
    email    text UNIQUE NOT NULL,

    CONSTRAINT right_nickname CHECK ( nickname ~* '[0-9a-z_]')
);

CREATE TABLE IF NOT EXISTS forum
(
    id      serial PRIMARY KEY,
    title   text        NOT NULL,
    profile text        NOT NULL,
    slug    text UNIQUE NOT NULL,
    -- posts
    -- threads

    FOREIGN KEY (profile) REFERENCES profile (nickname),
    CONSTRAINT good_slug CHECK (slug SIMILAR TO '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$')
);

CREATE TABLE IF NOT EXISTS thread
(
    id      SERIAL PRIMARY KEY,
    title   text NOT NULL,
    author  text NOT NULL,
    forum   text,
    message text NOT NULL,
    votes   int,
    slug    text,
    created timestamptz,

    FOREIGN KEY (author) REFERENCES profile (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug),
    CONSTRAINT right_slug CHECK ( slug SIMILAR TO '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$')
);

CREATE TABLE IF NOT EXISTS post
(
    id       serial PRIMARY KEY,
    parent   int           DEFAULT 0,
    author   text NOT NULL,
    message  text NOT NULL,
    isEdited bool NOT NULL DEFAULT false,
    forum    text,
    thread   int,
    created  timestamptz,

    FOREIGN KEY (parent) REFERENCES post (id),
    FOREIGN KEY (author) REFERENCES profile (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug),
    FOREIGN KEY (thread) REFERENCES thread (id)
);
