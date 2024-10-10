CREATE TABLE messages (
    id UUID PRIMARY KEY,         
    room_id UUID NOT NULL,                                          -- ID of the room where the message is sent (foreign key)
    user_id UUID NOT NULL,                                          -- ID of the user who sent the message (foreign key)
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT NOW(),                               -- Message content
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,   -- Reference to the rooms table
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE    -- Reference to the users table
);