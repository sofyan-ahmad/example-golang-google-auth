-- name: select-login
SELECT id, sub, givenName, familyName, profile, picture, email, emailVerified, gender, address, phone
    FROM users WHERE email = ? AND password = ? ;

-- name: select-email
SELECT id, sub, givenName, familyName, profile, picture, email, emailVerified, gender, address, phone
    FROM users WHERE email = ?;

-- name: insert
INSERT INTO users (id, sub, givenName, familyName, profile, picture, email, password, emailVerified, gender, address, phone) 
	values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: update
UPDATE users
    SET sub=?, givenName=?, familyName=?, profile=?, picture=?, email=?, emailVerified=?, gender=?, address=?, phone=?
    WHERE id=? 