-- name: select-login
SELECT id, sub, username, givenName, familyName, users.profile, picture, email, emailVerified, gender 
    FROM users WHERE email = ? AND password = ? ;

-- name: select-email
SELECT id, sub, username, givenName, familyName, users.profile, picture, email, emailVerified, gender 
    FROM users WHERE email = ?;

-- name: insert
INSERT INTO users (id, sub, username, givenName, familyName, users.profile, picture, email, password, emailVerified, gender) 
	values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);