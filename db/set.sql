CREATE EXTENSION IF NOT EXISTS CITEXT; -- eliminate calls to lower

CREATE UNLOGGED TABLE users
(
    Nickname   CITEXT PRIMARY KEY,
    FullName   TEXT NOT NULL,
    About      TEXT NOT NULL DEFAULT '',
    Email      CITEXT UNIQUE
);

CREATE UNLOGGED TABLE forum
(
    Title    TEXT   NOT NULL,
    "user"   CITEXT,
    Slug     CITEXT PRIMARY KEY,
    Posts    INT    DEFAULT 0,
    Threads  INT    DEFAULT 0
);

CREATE UNLOGGED TABLE threads
(
    Id      SERIAL    PRIMARY KEY,
    Title   TEXT      NOT NULL,
    Author  CITEXT    REFERENCES "users"(Nickname),
    Forum   CITEXT    REFERENCES "forum"(Slug),
    Message TEXT      NOT NULL,
    Votes   INT       DEFAULT 0,
    Slug    CITEXT,
    Created TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE UNLOGGED TABLE posts
(
    Id        SERIAL      PRIMARY KEY,
    Author    CITEXT,
    Created   TIMESTAMP   WITH TIME ZONE DEFAULT now(),
    Forum     CITEXT,
    IsEdited  BOOLEAN     DEFAULT FALSE,
    Message   CITEXT      NOT NULL,
    Parent    INT         DEFAULT 0,
    Thread    INT,
    Path      INTEGER[],
    FOREIGN KEY (thread) REFERENCES "threads" (id),
    FOREIGN KEY (author) REFERENCES "users"  (nickname)
);

CREATE UNLOGGED TABLE votes
(
    ID       SERIAL PRIMARY KEY,
    Author   CITEXT    REFERENCES "users" (Nickname),
    Voice    INT       NOT NULL,
    Thread   INT,
    FOREIGN KEY (thread) REFERENCES "threads" (id),
    UNIQUE (Author, Thread)
);


CREATE UNLOGGED TABLE users_forum
(
    Nickname  CITEXT  NOT NULL,
    FullName  TEXT    NOT NULL,
    About     TEXT,
    Email     CITEXT,
    Slug      CITEXT  NOT NULL,
    FOREIGN KEY (Nickname) REFERENCES "users" (Nickname),
    FOREIGN KEY (Slug) REFERENCES "forum" (Slug),
    UNIQUE (Nickname, Slug)
);

CREATE INDEX IF NOT EXISTS user_nickname ON users USING hash(nickname);
CREATE INDEX IF NOT EXISTS user_email ON users USING hash(email);
CREATE INDEX IF NOT EXISTS forum_slug ON forum USING hash(slug);
CREATE INDEX IF NOT EXISTS thr_date ON threads (created);
CREATE INDEX IF NOT EXISTS thr_forum_date ON threads(forum, created);
CREATE INDEX IF NOT EXISTS thr_forum ON threads USING hash(forum);
CREATE INDEX IF NOT EXISTS thr_slug ON threads USING hash(slug);
CREATE INDEX IF NOT EXISTS post_id_path ON posts(id, (path[1]));
CREATE INDEX IF NOT EXISTS post_thread_path_id ON posts(thread, path, id);
CREATE INDEX IF NOT EXISTS post_thread_id_path1_parent ON posts(thread, id, (path[1]), parent);
CREATE INDEX IF NOT EXISTS post_path1 ON posts((path[1]));
CREATE INDEX IF NOT EXISTS post_thr_id ON posts(thread);
CREATE INDEX IF NOT EXISTS post_thread_id ON posts(thread, id);
CREATE UNIQUE INDEX IF NOT EXISTS forum_users_unique ON users_forum(slug, nickname);
CREATE UNIQUE INDEX IF NOT EXISTS vote_unique ON votes(Author, Thread);

CREATE OR REPLACE FUNCTION insertVotes() RETURNS TRIGGER AS
$update_vote$
BEGIN
    UPDATE threads SET votes=(votes+NEW.voice) WHERE id=NEW.thread;
    return NEW;
end
$update_vote$ LANGUAGE plpgsql;

CREATE TRIGGER a_voice
    BEFORE INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insertVotes();


CREATE OR REPLACE FUNCTION updatePostUserForum() RETURNS TRIGGER AS
$update_forum_posts$
DECLARE
    t_fullname CITEXT;
    t_about    CITEXT;
    t_email CITEXT;
BEGIN
    SELECT fullname, about, email FROM users WHERE nickname = NEW.author INTO t_fullname, t_about, t_email;
    INSERT INTO users_forum (nickname, fullname, about, email, Slug)
    VALUES (New.Author, t_fullname, t_about, t_email, NEW.forum) on conflict do nothing;
    return NEW;
end
$update_forum_posts$ LANGUAGE plpgsql;

CREATE TRIGGER p_i_user_forum
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE updatePostUserForum();

CREATE OR REPLACE FUNCTION updateThreadUserForum() RETURNS TRIGGER AS
$update_forum_threads$
DECLARE
    a_nick CITEXT;
    t_fullname CITEXT;
    t_about    CITEXT;
    t_email CITEXT;
BEGIN
    SELECT Nickname, fullname, about, email
    FROM users WHERE Nickname = new.Author INTO a_nick, t_fullname, t_about, t_email;
    INSERT INTO users_forum (nickname, fullname, about, email, slug)
    VALUES (a_nick, t_fullname, t_about, t_email, NEW.forum) on conflict do nothing;
    return NEW;
end
$update_forum_threads$ LANGUAGE plpgsql;

CREATE TRIGGER t_i_forum_users
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE updateThreadUserForum();


CREATE OR REPLACE FUNCTION updateVotes() RETURNS TRIGGER AS
$update_votes$
BEGIN
    IF OLD.Voice <> NEW.Voice THEN
        UPDATE threads SET votes=(votes+NEW.Voice*2) WHERE id=NEW.Thread;
    END IF;
    return NEW;
end
$update_votes$ LANGUAGE plpgsql;

CREATE TRIGGER e_voice
    BEFORE UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE updateVotes();


CREATE OR REPLACE FUNCTION updateCountOfThreads() RETURNS TRIGGER AS
$update_forums$
BEGIN
    UPDATE forum SET Threads=(Threads+1) WHERE slug=NEW.forum;
    return NEW;
end
$update_forums$ LANGUAGE plpgsql;

CREATE TRIGGER a_t_i_forum
    BEFORE INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE updateCountOfThreads();


CREATE OR REPLACE FUNCTION updatePath() RETURNS TRIGGER AS
$update_path$
DECLARE
    parent_path  INTEGER[];
    parent_thread int;
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.path := array_append(new.path, new.id);
    ELSE
        SELECT thread FROM posts WHERE id = new.parent INTO parent_thread;
        IF NOT FOUND OR parent_thread != NEW.thread THEN
            RAISE EXCEPTION 'NOT FOUND OR parent_thread != NEW.thread' USING ERRCODE = '22409';
        end if;

        SELECT path FROM posts WHERE id = new.parent INTO parent_path;
        NEW.path := parent_path || new.id;
    END IF;
    UPDATE forum SET Posts=Posts + 1 WHERE forum.slug = new.forum;
    RETURN new;
END
$update_path$ LANGUAGE plpgsql;

CREATE TRIGGER u_p_trigger
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE updatePath();

CLUSTER users_forum USING forum_users_unique;

ANALYSE users_forum;
ANALYSE threads;
ANALYSE posts;