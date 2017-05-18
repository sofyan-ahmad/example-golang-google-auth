
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE users (
	id varchar(36) primary key,
	sub varchar(255),
	username varchar(255),
    givenName varchar(255), 
    familyName varchar(255),
    profile varchar(255), 
    picture  varchar(255), 
    email  varchar(255), 
    emailVerified  varchar(255), 
    gender   varchar(255),
	unique (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE users;