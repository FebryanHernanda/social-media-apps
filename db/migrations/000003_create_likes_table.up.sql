CREATE TABLE
    likes (
        id serial4 NOT NULL,
        post_id int4 NOT NULL,
        user_id int4 NOT NULL,
        liked_at timestamp DEFAULT now () NULL,
        CONSTRAINT likes_pkey PRIMARY KEY (id),
        CONSTRAINT unique_like UNIQUE (user_id, post_id),
        CONSTRAINT fk_likes_post FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
        CONSTRAINT fk_likes_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );