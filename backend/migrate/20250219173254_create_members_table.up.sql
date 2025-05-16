CREATE TABLE IF NOT EXISTS members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER NOT NULL,
		user TEXT NOT NULL,
		-- role TEXT NOT NULL, -- "member", "admin", "owner"
		UNIQUE(group_id, user)
	);