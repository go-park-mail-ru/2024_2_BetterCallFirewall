CREATE TABLE IF NOT EXISTS profile (
                                       id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                       first_name TEXT NOT NULL CONSTRAINT first_name_length CHECK (CHAR_LENGTH(first_name) <= 30),
                                       last_name TEXT NOT NULL CONSTRAINT last_name_length CHECK (CHAR_LENGTH(last_name) <= 30),
                                       email TEXT NOT NULL UNIQUE NOT NULL CONSTRAINT email_length CHECK (CHAR_LENGTH(email) <= 50),
                                       hashed_password TEXT NOT NULL,
                                       bio TEXT CONSTRAINT bio_length CHECK (CHAR_LENGTH(bio) <= 255) DEFAULT 'Что расскажете о себе?',
                                       avatar TEXT CONSTRAINT avatar_profile_length CHECK (CHAR_LENGTH(avatar) <= 100) DEFAULT '/image/default',
                                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS friend (
                                      sender INT,
                                      receiver INT CONSTRAINT unique_friends CHECK(sender != receiver),
                                      status INT,
                                      created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                      updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                      PRIMARY KEY (sender, receiver)
);

CREATE TABLE IF NOT EXISTS community (
                                         id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                         name text CONSTRAINT community_name_length CHECK(CHAR_LENGTH(name) <= 50),
                                         avatar TEXT CONSTRAINT avatar_community_length CHECK (CHAR_LENGTH(avatar) <= 100) DEFAULT '/image/default_community',
                                         about TEXT CONSTRAINT about_length CHECK (CHAR_LENGTH(about) <= 500) DEFAULT 'Опишите ваше сообщество',
                                         created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                         updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS admin (
    admin_id INT REFERENCES profile(id) ON DELETE CASCADE,
    community_id INT REFERENCES community(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS community_profile (
                                                 community_id INT REFERENCES community(id) ON DELETE CASCADE ,
                                                 profile_id INT REFERENCES profile(id) ON DELETE CASCADE ,
                                                 created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                                 updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS post (
                                    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                    author_id INT REFERENCES profile(id) ON DELETE CASCADE,
                                    community_id INT REFERENCES community(id) ON DELETE CASCADE DEFAULT NULL,
                                    content TEXT CONSTRAINT content_post_length CHECK (CHAR_LENGTH(content) <= 500) DEFAULT '',
                                    file_path TEXT CONSTRAINT file_path_length CHECK (CHAR_LENGTH(file_path) <= 100) DEFAULT '',
                                    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS message (
                                       id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                       receiver INT REFERENCES profile(id) ON DELETE CASCADE ,
                                       sender INT REFERENCES profile(id) ON DELETE CASCADE ,
                                       content TEXT CONSTRAINT content_message_length CHECK (CHAR_LENGTH(content) <= 500) DEFAULT '',
                                       file_path TEXT CONSTRAINT file_path_message_length CHECK (CHAR_LENGTH(file_path) <= 200) DEFAULT '',
                                       sticker_path TEXT CONSTRAINT sticker_path_message_length CHECK (CHAR_LENGTH(sticker_path) <= 200) DEFAULT '',
                                       is_read BOOLEAN DEFAULT FALSE,
                                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comment (
                                       id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                       user_id INT REFERENCES profile(id) ON DELETE CASCADE ,
                                       post_id INT REFERENCES post(id) ON DELETE CASCADE ,
                                       content TEXT CONSTRAINT content_comment_length CHECK (CHAR_LENGTH(content) <= 500) DEFAULT '',
                                       file_path TEXT CONSTRAINT file_path_comment CHECK ( CHAR_LENGTH(file_path) <= 200 ) DEFAULT '',
                                       sticker_path TEXT CONSTRAINT sticker_path_comment_length CHECK (CHAR_LENGTH(sticker_path) <= 200) DEFAULT '',
                                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reaction (
                                        post_id INT REFERENCES post(id) ON DELETE CASCADE DEFAULT NULL,
                                        comment_id INT REFERENCES comment(id) ON DELETE CASCADE DEFAULT NULL,
                                        user_id INT REFERENCES profile(id) ON DELETE CASCADE DEFAULT NULL,
                                        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE friend
    ADD FOREIGN KEY ("sender") REFERENCES profile(id) ON DELETE CASCADE,
    ADD FOREIGN KEY ("receiver") REFERENCES profile(id) ON DELETE CASCADE ;

ALTER TABLE reaction
    ADD CONSTRAINT unique_user_post UNIQUE (post_id, user_id),
    ADD CONSTRAINT unique_user_comment UNIQUE (comment_id, user_id);