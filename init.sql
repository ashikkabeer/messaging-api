CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    senderID VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
    receiverID VARCHAR(36) REFERENCES users(id)     ON DELETE CASCADE,
    content TEXT NOT NULL,
    read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- CREATE INDEX idx_messages_sender ON messages(senderID);
-- CREATE INDEX idx_messages_receiver ON messages(receiverID);

CREATE INDEX idx_messages_sender_receiver ON messages (senderID, receiverID);

INSERT INTO users (id) VALUES ('user123');
INSERT INTO users (id) VALUES ('user456');
INSERT INTO users (id) VALUES ('user789');
INSERT INTO users (id) VALUES ('user101');