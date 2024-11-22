-- users テーブル
CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    bio TEXT,
    profile_img_url VARCHAR(2083) -- URLの最大長に合わせる
);

-- posts テーブル
CREATE TABLE IF NOT EXISTS posts (
    post_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255),
    content VARCHAR(255),
    img_url VARCHAR(2083),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_at DATETIME,
    deleted_at DATETIME,
    parent_post_id VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_post_id) REFERENCES posts(post_id) ON DELETE SET NULL
);

-- likes テーブル
CREATE TABLE IF NOT EXISTS likes (
    user_id VARCHAR(255),
    post_id VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);

-- followers テーブル
CREATE TABLE IF NOT EXISTS followers (
    user_id VARCHAR(255),
    following_user_id VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, following_user_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (following_user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
