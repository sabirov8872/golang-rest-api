UPDATE users
SET firstname = $1,
    lastname = $2,
    username = $3,
    password = $4
WHERE id = $5