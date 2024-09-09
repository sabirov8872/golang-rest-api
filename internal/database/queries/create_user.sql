INSERT INTO users (firstname,
                   lastname,
                   username,
                   password)
VALUES ($1, $2, $3, $4)
RETURNING id