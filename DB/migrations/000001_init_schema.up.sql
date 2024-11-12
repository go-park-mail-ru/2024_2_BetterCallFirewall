CREATE TABLE IF NOT EXISTS profile (
                                       id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                       first_name TEXT NOT NULL CONSTRAINT first_name_length CHECK (CHAR_LENGTH(first_name) <= 30),
                                       last_name TEXT NOT NULL CONSTRAINT last_name_length CHECK (CHAR_LENGTH(last_name) <= 30),
                                       email TEXT NOT NULL UNIQUE NOT NULL CONSTRAINT email_length CHECK (CHAR_LENGTH(email) <= 50),
                                       hashed_password TEXT NOT NULL,
                                       bio TEXT CONSTRAINT bio_length CHECK (CHAR_LENGTH(bio) <= 255) DEFAULT 'Что расскажете о себе?',
                                       avatar INT DEFAULT 1,
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
                                         avatar INT DEFAULT '/default_community',
                                         about TEXT CONSTRAINT about_length CHECK (CHAR_LENGTH(about) <= 500) DEFAULT 'Опишите ваше сообщество',
                                         created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                         updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS community_profile (
                                                 community_id INT REFERENCES community(id) ON DELETE CASCADE ,
                                                 profile_id INT REFERENCES profile(id) ON DELETE CASCADE ,
                                                 created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                                 updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS file (
                                    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                    post_id INT DEFAULT NULL,
                                    comment_id INT DEFAULT NULL,
                                    profile_id INT DEFAULT NULL,
                                    file_path TEXT DEFAULT '/image/default',
                                    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS post (
                                    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                    author_id INT REFERENCES profile(id) ON DELETE CASCADE,
                                    community_id INT REFERENCES community(id) ON DELETE CASCADE DEFAULT NULL,
                                    content TEXT CONSTRAINT text_length CHECK (CHAR_LENGTH(content) <= 500) DEFAULT '',
                                    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS message (
                                    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                    receiver INT REFERENCES profile(id) ON DELETE CASCADE ,
                                    sender INT REFERENCES profile(id) ON DELETE CASCADE ,
                                    content TEXT CONSTRAINT content_length CHECK (CHAR_LENGTH(content) <= 500) DEFAULT '',
                                    is_read BOOLEAN DEFAULT FALSE,
                                    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comment (
                                       id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                       user_id INT REFERENCES profile(id) ON DELETE CASCADE ,
                                       post_id INT REFERENCES post(id) ON DELETE CASCADE ,
                                       content TEXT CONSTRAINT content_length CHECK (CHAR_LENGTH(content) <= 500) DEFAULT '',
                                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reaction (
                                        post_id INT REFERENCES post(id) ON DELETE CASCADE DEFAULT NULL,
                                        comment_id INT REFERENCES comment(id) ON DELETE CASCADE DEFAULT NULL,
                                        user_id INT REFERENCES profile(id) ON DELETE CASCADE DEFAULT NULL,
                                        file_id INT REFERENCES file(id) ON DELETE CASCADE DEFAULT NULL,
                                        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE profile
    ADD FOREIGN KEY ("avatar") REFERENCES file(id) ON DELETE SET NULL;

ALTER TABLE friend
    ADD FOREIGN KEY ("sender") REFERENCES profile(id) ON DELETE CASCADE,
    ADD FOREIGN KEY ("receiver") REFERENCES profile(id) ON DELETE CASCADE ;

ALTER TABLE community ADD FOREIGN KEY ("avatar") REFERENCES file(id) ON DELETE SET NULL;

ALTER TABLE file
    ADD FOREIGN KEY ("post_id") REFERENCES post(id) ON DELETE NO ACTION,
    ADD FOREIGN KEY ("comment_id") REFERENCES comment(id) ON DELETE NO ACTION,
    ADD FOREIGN KEY ("profile_id") REFERENCES profile(id) ON DELETE NO ACTION;

INSERT INTO file(file_path) VALUES ('/image/default');