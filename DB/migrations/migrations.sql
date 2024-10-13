CREATE TABLE IF NOT EXISTS profile (
    id INT PRIMARY KEY,
    first_name TEXT NOT NULL CONSTRAINT first_name_length CHECK (CHAR_LENGTH(first_name) <= 30),
    last_name TEXT NOT NULL CONSTRAINT last_name_length CHECK (CHAR_LENGTH(last_name) <= 30),
    email TEXT NOT NULL UNIQUE NOT NULL CONSTRAINT email_length CHECK (CHAR_LENGTH(email) <= 50),
    hashed_password TEXT NOT NULL,
    bio TEXT CONSTRAINT bio_length CHECK (CHAR_LENGTH(bio) <= 255) DEFAULT 'description',
    avatar INT REFERENCES file(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS friend (
    sender INT REFERENCES profile(id),
    receiver INT REFERENCES profile(id) CONSTRAINT unique_friends CHECK(sender != receiver),
    status INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS community (
    id INT PRIMARY KEY,
    name text CONSTRAINT community_name_length CHECK(CHAR_LENGTH(name) <= 50),
    avatar INT REFERENCES file(id),
    about TEXT CONSTRAINT about_length CHECK (CHAR_LENGTH(about) <= 500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS community_profile (
    community_id INT REFERENCES community(id),
    profile_id INT REFERENCES profile(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS file (
    id INT PRIMARY KEY,
    post_id INT,
    comment_id INT,
    profile_id INT REFERENCES profile(id),
    likes INT DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS post (
    id INT PRIMARY KEY,
    author_id INT REFERENCES profile(id),
    community_id INT REFERENCES community(id),
    content TEXT CONSTRAINT text_length CHECK (CHAR_LENGTH(content) <= 500),
    likes INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS message (
    receiver INT REFERENCES profile(id),
    sender INT REFERENCES profile(id),
    content TEXT CONSTRAINT content_length CHECK (CHAR_LENGTH(content) <= 500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comment (
    id INT PRIMARY KEY,
    user_id INT REFERENCES profile(id),
    post_id INT REFERENCES post(id),
    likes INT DEFAULT 0,
    content TEXT CONSTRAINT content_length CHECK (CHAR_LENGTH(content) <= 500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reaction (
    post_id INT REFERENCES post(id),
    comment_id INT REFERENCES comment(id),
    user_id INT REFERENCES profile(id),
    file_id INT REFERENCES file(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
)

