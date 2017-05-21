
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users
ADD COLUMN phone varchar(255),
ADD COLUMN address varchar(255);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE users
DROP COLUMN phone varchar(255),
DROP COLUMN password;