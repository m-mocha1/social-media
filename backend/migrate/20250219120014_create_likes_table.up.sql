CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		com_id INTEGER,
		username TEXT NOT NULL,
		UNIQUE(com_id,post_id,username)
	);