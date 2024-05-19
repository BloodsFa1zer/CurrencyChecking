CREATE TABLE IF NOT EXISTS Users (
    ID SERIAL primary key,
    email TEXT UNIQUE not null,
    created_at text not null
);