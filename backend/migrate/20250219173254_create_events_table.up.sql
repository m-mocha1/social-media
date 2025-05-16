CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    creator TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    event_datetime DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (creator) REFERENCES users(username) ON DELETE CASCADE
);