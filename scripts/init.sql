DROP TABLE IF EXISTS users, forums, posts, threads CASCADE;

CREATE TABLE IF NOT EXISTS users
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

CREATE TABLE IF NOT EXISTS forums
(
    id      serial PRIMARY KEY,
    title   text        NOT NULL,
    profile text        NOT NULL,
    slug    text UNIQUE NOT NULL,
    -- posts
    -- threads

    FOREIGN KEY (profile) REFERENCES users (nickname)
    -- CONSTRAINT good_slug CHECK (slug SIMILAR TO '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$')
);

CREATE TABLE IF NOT EXISTS threads
(
    id      SERIAL PRIMARY KEY,
    title   text NOT NULL,
    author  text NOT NULL,
    forum   text,
    message text NOT NULL,
    votes   int,
    slug    text,
    created timestamptz,

    FOREIGN KEY (author) REFERENCES users (nickname),
    FOREIGN KEY (forum) REFERENCES forums (slug)
    -- CONSTRAINT right_slug CHECK ( slug SIMILAR TO '^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$')
);

CREATE TABLE IF NOT EXISTS votes
(
    thread_id int  NOT NULL,
    user_id   int  NOT NULL,
    likes     bool NOT NULL,

    PRIMARY KEY (thread_id, user_id),
    FOREIGN KEY (thread_id) REFERENCES threads (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS posts
(
    id       serial PRIMARY KEY,
    parent   int           DEFAULT 0,
    author   text NOT NULL,
    message  text NOT NULL,
    isEdited bool NOT NULL DEFAULT false,
    forum    text,
    thread   int,
    created  timestamptz,

    UNIQUE (id, parent, author, thread),
    FOREIGN KEY (author) REFERENCES users (nickname),
    FOREIGN KEY (forum) REFERENCES forums (slug),
    FOREIGN KEY (thread) REFERENCES threads (id)
);


CREATE OR REPLACE FUNCTION votes_ins_upd() RETURNS trigger AS
$$
DECLARE
    value int;
BEGIN
    IF TG_OP = 'UPDATE' THEN
        value := 2;
    ELSE
        value := 1;
    END IF;

    IF NEW.likes = TRUE THEN
        UPDATE threads
        SET votes = votes + value
        WHERE id = NEW.thread_id;
    ELSE
        UPDATE threads
        SET votes = votes - value
        WHERE id = NEW.thread_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION votes_del() RETURNS trigger AS
$$
DECLARE
    var_likes boolean;
BEGIN
    SELECT likes
    FROM votes
    WHERE thread_id = OLD.thread_id
      AND user_id = OLD.user_id
    INTO STRICT var_likes;

    IF var_likes = TRUE THEN
        UPDATE threads
        SET votes = votes - 1
        WHERE id = OLD.thread_id;
    ELSE
        UPDATE threads
        SET votes = votes + 1
        WHERE id = OLD.thread_Id;
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER votes_ins_upd AFTER INSERT OR UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE votes_ins_upd();
CREATE TRIGGER vots_del BEFORE DELETE ON votes
    FOR EACH ROW EXECUTE PROCEDURE votes_del();
