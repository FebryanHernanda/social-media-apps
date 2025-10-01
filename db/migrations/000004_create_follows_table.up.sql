CREATE TABLE
    follows (
        id serial4 NOT NULL,
        user_id int4 NOT NULL,
        followed_user_id int4 NOT NULL,
        created_at timestamp DEFAULT now () NULL,
        CONSTRAINT follows_pkey PRIMARY KEY (id),
        CONSTRAINT no_self_follow CHECK ((user_id <> followed_user_id)),
        CONSTRAINT unique_follow UNIQUE (user_id, followed_user_id),
        CONSTRAINT fk_follow_followed FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
        CONSTRAINT fk_follow_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );