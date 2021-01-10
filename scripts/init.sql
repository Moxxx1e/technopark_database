CREATE EXTENSION IF NOT EXISTS citext;
DROP TABLE IF EXISTS users, forums, posts, threads, votes, user_forum CASCADE;

CREATE UNLOGGED TABLE IF NOT EXISTS users
(
    id       serial PRIMARY KEY,
    nickname citext COLLATE "POSIX" UNIQUE NOT NULL,
    fullname text                          NOT NULL,
    about    text                          NOT NULL,
    email    citext UNIQUE                 NOT NULL
);

CREATE INDEX users_cover ON users (nickname, fullname, about, email);
CREATE INDEX users_nickname ON users using hash (nickname);
CREATE INDEX users_email ON users using hash (email);

CREATE UNLOGGED TABLE IF NOT EXISTS forums
(
    id      serial PRIMARY KEY,
    title   text          NOT NULL,
    profile citext        NOT NULL,
    slug    citext UNIQUE NOT NULL,
    posts   int DEFAULT 0,
    threads int DEFAULT 0,

    FOREIGN KEY (profile) REFERENCES users (nickname)
);
CREATE INDEX forums_cover ON forums (title, profile, slug, posts, threads);
CREATE INDEX forums_slug ON forums USING hash (slug);
CREATE INDEX forums_user ON forums (profile);

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
CREATE INDEX threads_cover ON threads (title, author, forum, message, votes, slug, created);
CREATE INDEX threads_created_forum ON threads (created, forum);
CREATE INDEX threads_created ON threads (created);
CREATE INDEX threads_slug ON threads using hash (slug);
CREATE INDEX threads_author ON threads (author);
CREATE INDEX threads_forum ON threads (forum);

CREATE UNLOGGED TABLE IF NOT EXISTS votes
(
    thread_id int  NOT NULL,
    user_id   int  NOT NULL,
    likes     bool NOT NULL,

    PRIMARY KEY (thread_id, user_id),
    FOREIGN KEY (thread_id) REFERENCES threads (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE UNIQUE INDEX votes_thread_user ON votes (thread_id, user_id);

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
CREATE INDEX posts_thread_id on posts (thread, id);
CREATE INDEX posts_path on posts (path);
CREATE INDEX posts_path1_path on posts ((path[1]), path);
CREATE INDEX posts_created ON posts (created);
CREATE INDEX IF NOT EXISTS posts_cover
    ON posts (id, parent, path, author, message, isEdited, forum, thread, created);
CREATE INDEX posts_thread ON posts (thread);
CREATE INDEX posts_forum ON posts (forum);

CREATE UNLOGGED TABLE IF NOT EXISTS user_forum
(
    nickname citext,
    slug     citext,

    PRIMARY KEY (nickname, slug),
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    FOREIGN KEY (slug) REFERENCES forums (slug)
);
CREATE INDEX user_forum_nickname ON user_forum (nickname);
CREATE INDEX user_forum_slug ON user_forum (slug);
CREATE INDEX user_forum_nickname_slug ON user_forum (nickname, slug);

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
        RAISE EXCEPTION 'Parent post does not exist in thread';
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

CREATE OR REPLACE FUNCTION user_forum_ins() RETURNS trigger AS
$ins_author$
BEGIN
    INSERT INTO user_forum(nickname, slug)
    VALUES (NEW.author, NEW.forum)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$ins_author$
    LANGUAGE plpgsql;

CREATE TRIGGER user_forum_ins_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE user_forum_ins();

CREATE TRIGGER user_forum_ins_posts
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE user_forum_ins();

CREATE OR REPLACE FUNCTION posts_inc() RETURNS trigger AS
$$
BEGIN
    UPDATE forums
    SET posts = posts + 1
    WHERE slug = NEW.forum;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER posts_inc
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE posts_inc();

CREATE OR REPLACE FUNCTION threads_inc() RETURNS trigger AS
$$
BEGIN
    UPDATE forums
    SET threads = threads + 1
    WHERE slug = NEW.forum;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER threads_inc
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE threads_inc();
