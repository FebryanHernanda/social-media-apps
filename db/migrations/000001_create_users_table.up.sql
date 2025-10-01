CREATE TABLE
    users (
        id serial4 NOT NULL,
        email varchar(255) NOT NULL,
        "password" varchar(255) NOT NULL,
        "name" varchar(255) NOT NULL,
        avatar_path text NULL,
        created_at timestamp DEFAULT now () NULL,
        deleted_at timestamp NULL,
        biography text NULL,
        CONSTRAINT users_pkey PRIMARY KEY (id)
    );