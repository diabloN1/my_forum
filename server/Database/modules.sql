-- 1. Users 👤​
CREATE TABLE IF NOT EXISTS users (
    id TEXT UNIQUE PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    user_name TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    user_image TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Posts 📝​
CREATE TABLE IF NOT EXISTS posts (
    id TEXT UNIQUE PRIMARY KEY,
    user_id TEXT NOT NULL,
    title TEXT,
    content TEXT,
    image_url TEXT, -- URL or path to the post image
    category TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 3. Comments 💭​
CREATE TABLE IF NOT EXISTS comments (
    id TEXT UNIQUE PRIMARY KEY,
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 4. Categories 🏷️​
CREATE TABLE IF NOT EXISTS categories (
    id TEXT UNIQUE PRIMARY KEY,
    category_name TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 5. PostCategory 🔗​
CREATE TABLE IF NOT EXISTS postCategory (
    post_id TEXT NOT NULL,
    category_id TEXT NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- 6. LikeDislike 👍👎​
CREATE TABLE IF NOT EXISTS likeDislike (
    id TEXT UNIQUE PRIMARY KEY, 
    user_id TEXT NOT NULL,
    post_id TEXT NOT NULL,
    is_like BOOLEAN NOT NULL, 
    is_comment BOOLEAN NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);


-- 7. Session 🍪​
CREATE TABLE IF NOT EXISTS Session (
    id TEXT UNIQUE PRIMARY KEY,
    user_id TEXT NOT NULL,
    user_ip TEXT UNIQUE NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- 8. Conversations 💬​
CREATE TABLE IF NOT EXISTS Conversations (
    id TEXT UNIQUE PRIMARY KEY,
    user1_id TEXT NOT NULL,
    user2_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user1_id) REFERENCES Users(id),
    FOREIGN KEY (user2_id) REFERENCES Users(id)
);

-- 9. Messages ✉️​
CREATE TABLE IF NOT EXISTS Messages (
    id TEXT UNIQUE PRIMARY KEY,
    conversation_id TEXT NOT NULL,
    sender_id TEXT NOT NULL,
    content TEXT NOT NULL,
    sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (conversation_id) REFERENCES Conversations(id),
    FOREIGN KEY (sender_id) REFERENCES Users(id)
);
