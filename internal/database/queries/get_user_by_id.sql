SELECT id,
       firstname,
       lastname,
       username,
       password
FROM users
WHERE id = $1