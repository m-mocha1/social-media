CREATE TABLE not_requests (
    id  INTEGER PRIMARY KEY AUTOINCREMENT,    
    actionId TEXT, 
    sender TEXT NOT NULL,      
    receiver TEXT NOT NULL,
    nottype TEXT NOT NULL, 
    status TEXT DEFAULT 'pending', 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(sender, receiver, nottype)
);
