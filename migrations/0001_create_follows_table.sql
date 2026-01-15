CREATE TABLE follows (
    follower_id INT NOT NULL,
    seller_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (follower_id, seller_id),
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREING KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE
);