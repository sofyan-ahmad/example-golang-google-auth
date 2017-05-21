-- name: select-login
SELECT id, sub, givenName, familyName, profile, picture, email, emailVerified, gender, address, phone
    FROM users WHERE email = ? AND password = ? ;

-- name: select-email
SELECT id, sub, givenName, familyName, profile, picture, email, emailVerified, gender, address, phone
    FROM users WHERE email = ?;

-- name: select-reset-token
SELECT id, sub, givenName, familyName, profile, picture, email, emailVerified, gender, address, phone
    FROM users WHERE email = ? AND reset_token = ? ;

-- name: insert
INSERT INTO users (id, sub, givenName, familyName, profile, picture, email, password, emailVerified, gender, address, phone) 
	values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: update
UPDATE users
    SET sub=?, givenName=?, familyName=?, profile=?, picture=?, email=?, emailVerified=?, gender=?, address=?, phone=?
    WHERE id=? 

-- name: update-reset-token
UPDATE users 
    SET reset_token = ? WHERE id=? 


-- name: update-password
UPDATE users 
    SET password = ?, reset_token = "" WHERE id=? 