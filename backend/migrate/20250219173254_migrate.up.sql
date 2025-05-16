CREATE TABLE IF NOT EXISTS follows(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_username TEXT,
    following_username TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     CONSTRAINT unique_follow UNIQUE ( follower_username,  following_username)
);
