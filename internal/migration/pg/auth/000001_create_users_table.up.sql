CREATE TABLE IF NOT EXISTS users
(
    username        TEXT PRIMARY KEY,
    password        TEXT NOT NULL,
    role            TEXT NOT NULL,
    plan_name       TEXT NOT NULL,
    links_remaining INT  NOT NULL
);

INSERT INTO users (username, password, role, plan_name, links_remaining)
VALUES ('admin', '$2a$10$cFmATgr3Sv9TRDHUQ64qZeEguzTMDIJKWz/euADw0dB5D2Lc5NjYm', 'admin', 'premium', 1000);