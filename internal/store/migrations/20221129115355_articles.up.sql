CREATE TABLE IF NOT EXISTS article(
   id serial PRIMARY KEY,
   date_at VARCHAR NOT NULL unique,
   title VARCHAR NOT NULL,
   explanation VARCHAR NOT NULL
);