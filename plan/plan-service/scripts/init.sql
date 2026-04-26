CREATE TABLE IF NOT EXISTS plan (
    id        VARCHAR(36)  PRIMARY KEY NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    title     VARCHAR(255) NOT NULL,
    date      TIMESTAMP
);

CREATE TABLE IF NOT EXISTS todo (
    id        VARCHAR(36)  PRIMARY KEY NOT NULL,
    plan_id   VARCHAR(36),
    user_name VARCHAR(255) NOT NULL,
    date      TIMESTAMP,
    content   TEXT,
    FOREIGN KEY (plan_id) REFERENCES plan(id) ON DELETE CASCADE
);
