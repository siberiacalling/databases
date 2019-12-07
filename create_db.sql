CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
  id       SERIAL NOT NULL
    CONSTRAINT users_pkey
    PRIMARY KEY,
  nickname CITEXT NOT NULL,
  fullname TEXT   NOT NULL,
  about    TEXT,
  email    CITEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS forum
(
  id      SERIAL NOT NULL
    CONSTRAINT forum_pkey
    PRIMARY KEY,
  title   CITEXT NOT NULL,
  author  CITEXT NOT NULL,
  slug    CITEXT NOT NULL,
  posts   BIGINT DEFAULT 0,
  threads BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS thread
(
  id      SERIAL NOT NULL
    CONSTRAINT thread_pkey
    PRIMARY KEY,
  title   CITEXT NOT NULL,
  author  CITEXT NOT NULL,
  forum   CITEXT NOT NULL,
  message CITEXT NOT NULL,
  votes   INTEGER                  DEFAULT 0,
  created TIMESTAMP WITH TIME ZONE DEFAULT now(),
  slug    CITEXT
);

CREATE TABLE IF NOT EXISTS post
(
  id        SERIAL            NOT NULL
    CONSTRAINT post_pkey
    PRIMARY KEY,
  parent    INTEGER DEFAULT 0 NOT NULL,
  author    CITEXT,
  message   CITEXT,
  is_edited BOOLEAN                  DEFAULT FALSE,
  forum     CITEXT            NOT NULL,
  thread    INTEGER                  DEFAULT 0,
  created   TIMESTAMP WITH TIME ZONE DEFAULT now(),
  path      INTEGER []               DEFAULT ARRAY [] :: INTEGER [],
  root      INTEGER                  DEFAULT 0
);

CREATE TABLE voice
(
  id         SERIAL            NOT NULL
    CONSTRAINT voice_pkey
    PRIMARY KEY,
  nickname   CITEXT,
  vote       INTEGER DEFAULT 0 NOT NULL,
  prev_vote  INTEGER                  DEFAULT 0,
  thread_id  INTEGER                  DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS user_nickname_uindex
  ON users (nickname);

CREATE UNIQUE INDEX IF NOT EXISTS user_email_uindex
  ON users (email);

CREATE UNIQUE INDEX IF NOT EXISTS forum_slug_uindex
  ON forum (slug);