CREATE TABLE IF NOT EXISTS users (
  id serial PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  full_name TEXT,
  password TEXT,
  business_name TEXT
);
