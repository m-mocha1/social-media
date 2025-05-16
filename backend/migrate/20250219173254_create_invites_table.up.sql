CREATE TABLE IF NOT EXISTS invites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER NOT NULL,
		inviter TEXT NOT NULL DEFAULT 'himself',
		invitee TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending' -- "pending", "accepted", "rejected"
	);