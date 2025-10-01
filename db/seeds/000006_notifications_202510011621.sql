INSERT INTO
	public.notifications (
		id,
		receiver_id,
		actor_id,
		action_type,
		post_id,
		is_read,
		created_at
	)
VALUES
	(
		2,
		2,
		1,
		'follow',
		NULL,
		true,
		'2025-10-01 02:26:04.884864'
	),
	(
		10,
		2,
		1,
		'comment',
		2,
		false,
		'2025-10-01 04:58:38.458089'
	),
	(
		11,
		1,
		2,
		'comment',
		1,
		false,
		'2025-10-01 05:00:45.653283'
	),
	(
		12,
		1,
		2,
		'comment',
		1,
		false,
		'2025-10-01 05:00:47.479198'
	),
	(
		15,
		2,
		1,
		'comment',
		2,
		false,
		'2025-10-01 06:33:57.834764'
	),
	(
		16,
		2,
		1,
		'comment',
		3,
		false,
		'2025-10-01 06:34:26.477151'
	),
	(
		17,
		2,
		1,
		'like',
		2,
		false,
		'2025-10-01 07:45:41.409663'
	),
	(
		13,
		1,
		2,
		'comment',
		1,
		true,
		'2025-10-01 05:00:48.979899'
	);