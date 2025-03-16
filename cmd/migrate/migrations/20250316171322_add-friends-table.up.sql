CREATE TABLE friends (
                         id SERIAL PRIMARY KEY,
                         user_id VARCHAR(15) NOT NULL,
                         friend_id VARCHAR(15) NOT NULL,
                         status TEXT NOT NULL CHECK (status IN ('pending', 'accepted', 'rejected', 'blocked')),
                         created_at TIMESTAMP DEFAULT NOW(),
                         UNIQUE (user_id, friend_id),
                         FOREIGN KEY (user_id) REFERENCES users(noxId) ON DELETE CASCADE,
                         FOREIGN KEY (friend_id) REFERENCES users(noxId) ON DELETE CASCADE
);