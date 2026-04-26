CREATE TABLE IF NOT EXISTS todo (
    id          VARCHAR(36)  PRIMARY KEY NOT NULL,
    user_name   VARCHAR(255) NOT NULL,
    date        DATE,
    content     TEXT
);
