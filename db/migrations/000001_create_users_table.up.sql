CREATE TABLE users (
    id UUID PRIMARY KEY,         
    username VARCHAR(255) NOT NULL UNIQUE, -- Unique username
    email VARCHAR(255) NOT NULL UNIQUE,    -- Unique email
    password VARCHAR(255) NOT NULL        
);
