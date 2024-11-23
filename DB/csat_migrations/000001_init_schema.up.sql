CREATE TABLE IF NOT EXISTS csat_metric (
                                           total INT DEFAULT NULL,
                                           review TEXT CONSTRAINT review_length CHECK (CHAR_LENGTH(review) <= 500) DEFAULT '',
                                           created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);