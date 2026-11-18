-- +goose Up
-- Существующим пользователям выдаём временный уникальный логин из их id,
-- чтобы можно было включить NOT NULL. Формат совпадает с правилами логина
-- (начинается с буквы, только [a-z0-9_]).
ALTER TABLE users ADD COLUMN login TEXT;
UPDATE users SET login = 'u_' || replace(id::text, '-', '') WHERE login IS NULL;
ALTER TABLE users ALTER COLUMN login SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT users_login_key UNIQUE (login);

-- +goose Down
ALTER TABLE users DROP CONSTRAINT users_login_key;
ALTER TABLE users DROP COLUMN login;
