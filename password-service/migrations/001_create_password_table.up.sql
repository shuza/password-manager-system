CREATE TABLE IF NOT EXISTS passwords (
    id serial PRIMARY KEY,
    user_id INT NOT NULL,
    account_name TEXT NOT NULL,
    password TEXT NOT NULL,
    email TEXT,
    username TEXT
);
