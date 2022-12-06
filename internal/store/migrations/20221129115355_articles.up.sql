CREATE TABLE IF NOT EXISTS article(
   date_at VARCHAR NOT NULL unique,
   title VARCHAR NOT NULL,
   explanation VARCHAR NOT NULL
);