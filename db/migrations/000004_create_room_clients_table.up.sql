CREATE TABLE room_clients (   
    room_id UUID NOT NULL,                                           -- ID of the room (foreign key)
    user_id UUID NOT NULL,                                           -- ID of the user (foreign key)
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,   -- Reference to the rooms table
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE    -- Reference to the users table
);