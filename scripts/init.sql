CREATE EXTENSION IF NOT EXISTS citext;
DROP TABLE IF EXISTS users, forums, posts, threads, votes, user_forum CASCADE;

CREATE UNLOGGED TABLE IF NOT EXISTS users
(
    id       serial PRIMARY KEY,
    nickname citext COLLATE "POSIX" UNIQUE NOT NULL,
    fullname text          NOT NULL,
    about    text          NOT NULL,
    email    citext UNIQUE NOT NULL
);

CREATE UNLOGGED TABLE IF NOT EXISTS forums
(
    id      serial PRIMARY KEY,
    title   text          NOT NULL,
    profile citext        NOT NULL,
    slug    citext UNIQUE NOT NULL,
    -- posts
    -- threads

    FOREIGN KEY (profile) REFERENCES users (nickname)
);

CREATE UNLOGGED TABLE IF NOT EXISTS threads
(
    id      SERIAL PRIMARY KEY,
    title   text   NOT NULL,
    author  citext NOT NULL,
    forum   citext,
    message text   NOT NULL,
    votes   int,
    slug    citext,
    created timestamptz,

    FOREIGN KEY (author) REFERENCES users (nickname),
    FOREIGN KEY (forum) REFERENCES forums (slug)
);

CREATE UNLOGGED TABLE IF NOT EXISTS votes
(
    thread_id int  NOT NULL,
    user_id   int  NOT NULL,
    likes     bool NOT NULL,

    PRIMARY KEY (thread_id, user_id),
    FOREIGN KEY (thread_id) REFERENCES threads (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE UNLOGGED TABLE IF NOT EXISTS posts
(
    id       serial PRIMARY KEY,
    parent   int    NOT NULL,
    path     int[]  NOT NULL,
    author   citext NOT NULL,
    message  text   NOT NULL,
    isEdited bool   NOT NULL DEFAULT false,
    forum    citext,
    thread   int,
    created  timestamptz,

    FOREIGN KEY (author) REFERENCES users (nickname),
    FOREIGN KEY (forum) REFERENCES forums (slug),
    FOREIGN KEY (thread) REFERENCES threads (id)
);

CREATE UNLOGGED TABLE IF NOT EXISTS user_forum
(
    nickname citext,
    slug     citext,

    PRIMARY KEY (nickname, slug),
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    FOREIGN KEY (slug) REFERENCES forums (slug)
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

CREATE TRIGGER votes_ins_upd
    AFTER INSERT OR UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE votes_ins_upd();
CREATE TRIGGER vots_del
    BEFORE DELETE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE votes_del();

-- Insert id into path on insert
CREATE OR REPLACE FUNCTION upd_path() RETURNS trigger AS
$upd_path$
DECLARE
    parent_thread integer;
    parent_path   integer[];
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.path := array_append(NEW.path, NEW.id);
        RETURN NEW;
    END IF;

    SELECT thread INTO parent_thread FROM posts WHERE id = NEW.parent;
    IF NOT FOUND OR NEW.thread <> parent_thread THEN
        RAISE EXCEPTION 'Can not find parent post into thread';
    END IF;

    SELECT path INTO parent_path FROM posts WHERE id = NEW.parent;
    NEW.path = array_append(parent_path, NEW.id);
    RETURN NEW;
END;
$upd_path$
    LANGUAGE plpgsql;

CREATE TRIGGER upd_path
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE upd_path();

CREATE OR REPLACE FUNCTION ins_author() RETURNS trigger AS
$ins_author$
BEGIN
    INSERT INTO user_forum(nickname, slug)
    VALUES(NEW.author, NEW.forum)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$ins_author$
    LANGUAGE plpgsql;

CREATE TRIGGER ins_author_on_ins_thread AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE ins_author();

CREATE TRIGGER ins_author_on_ins_post AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE ins_author();
