CREATE TABLE
    posts (
        id serial4 NOT NULL,
        "content" text NULL,
        image_path text NULL,
        user_id int4 NOT NULL,
        created_at timestamp DEFAULT now () NULL,
        deleted_at timestamp NULL,
        CONSTRAINT posts_pkey PRIMARY KEY (id),
        CONSTRAINT fk_post_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );