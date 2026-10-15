-- +goose Up
 CREATE Table users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT not NULL 
 );

-- +goose Down
DROP TABLE users;
