CREATE TABLE notifications (
	id serial4 NOT NULL,
	receiver_id int4 NOT NULL,
	actor_id int4 NOT NULL,
	action_type varchar(50) NOT NULL,
	post_id int4 NULL,
	is_read bool DEFAULT false NULL,
	created_at timestamp DEFAULT now() NULL,
	CONSTRAINT notifications_action_type_check CHECK (((action_type)::text = ANY ((ARRAY['like'::character varying, 'comment'::character varying, 'follow'::character varying])::text[]))),
	CONSTRAINT notifications_pkey PRIMARY KEY (id),
	CONSTRAINT fk_notif_actor FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT fk_notif_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
	CONSTRAINT fk_notif_receiver FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
);