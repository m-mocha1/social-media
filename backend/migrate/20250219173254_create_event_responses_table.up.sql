CREATE TABLE event_responses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    user TEXT NOT NULL,
    response TEXT CHECK(response IN ('Going', 'Not going')) NOT NULL,
    responded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (user) REFERENCES users(username) ON DELETE CASCADE,
    UNIQUE(event_id, user) -- Ensure a user can only respond once per event
);