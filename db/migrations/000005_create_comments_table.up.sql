CREATE TABLE
    "comments" (
        id serial4 NOT NULL,
        "content" text NOT NULL,
        post_id int4 NOT NULL,
        user_id int4 NOT NULL,
        created_at timestamp DEFAULT now () NULL,
        deleted_at timestamp NULL,
        CONSTRAINT comments_pkey PRIMARY KEY (id),
        CONSTRAINT fk_comment_post FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
        CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );