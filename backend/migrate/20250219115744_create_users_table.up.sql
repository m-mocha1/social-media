CREATE TABLE IF NOT EXISTS users (
       	id INTEGER PRIMARY KEY AUTOINCREMENT,
		pfp BLOB,
		username TEXT NOT NULL,
		age TEXT NOT NULL,
		gender TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,	
		public TEXT,
		aboutme TEXT,
		UNIQUE(username)
		UNIQUE(email)
);
