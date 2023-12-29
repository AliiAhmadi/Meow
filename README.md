# Meow

![golang](https://github.com/AliiAhmadi/Meow/assets/107758775/28b90dae-bc10-4376-ba30-938982142b72)

> [!TIP]
> You can use `make` to run the project.

Let's go through the general process starting from creating a new account and verifying the account and using it.

## Registeration

For registration we need to send a post request to `http://127.0.0.1:4000/v1/users` with `email`, `password` and `name` fields.

```zsh
curl --location 'http://127.0.0.1:4000/v1/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "me@gmail.com",
    "name": "me me",
    "password": "123456789"
}'
```

Response will be like this if registration be sucessful.
```json
{
    "user": {
        "id": 17,
        "created_at": "2023-12-26T13:16:56+03:30",
        "name": "me me",
        "email": "me@gmail.com",
        "activated": false
    }
}
```

Now you will receive an verification email like this:
```
Hi,

Thanks for signing up.

For future reference, your user ID number is 17.

Please send a request to the PUT /v1/users/activated endpoint with the following JSON body to activate your account:


{"token": "TXF7POG4374WTUUS2Y7XQ5ETFU"}

Please note that this is a one-time use token and it will expire in 1 hour.

Thanks
```

Follow this and send a PUT request to `/v1/users/activated`:

```zsh
curl --location --request PUT 'http://127.0.0.1:4000/v1/users/activated' \
--header 'Key: Value' \
--header 'Content-Type: application/json' \
--data '{
    "token": "TXF7POG4374WTUUS2Y7XQ5ETFU"
}'
```

Your account activated and now you need to generate a authorization token for access to endpoints:

```zsh
curl --location 'http://127.0.0.1:4000/v1/tokens/authentication' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "me@gmail.com",
    "password": "123456789"
}'
```

Response:

```json
{
    "auth_token": {
        "token": "TZQQKL3PYCIFROF7INLAYM5YPI",
        "expiry": "2023-12-27T13:29:45.791141968+03:30"
    }
}
```

Now you have a bearer token. use that for every request. If you send a request to an protected route without any authirization token you will get error:

```zsh
curl --location 'http://127.0.0.1:4000/v1/movies/1' \
--header 'Content-Type: application/json'
```

Response (with 401 status code):

```json
{
    "error": "you must be authenticated to access this resource"
}
```

Or if use an fake token:

```zsh
curl --location 'http://127.0.0.1:4000/v1/movies/1' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer 3PXHUW7RNJCOY5L7JIP7NXKZZ4'
```

```json
{
    "error": "invalid or missing authentication token"
}
```

Some validation will be done when you want to register:

```zsh
curl --location 'http://127.0.0.1:4000/v1/users' \
--header 'Content-Type: application/json' \
--data '{
    "email": "invalid_email",
    "password": "123_"
}'
```

```json
{
    "error": {
        "email": "must be a valid email address",
        "name": "must be provided",
        "password": "must be at least 8 bytes long"
    }
}
```

## Movie routes

We have 5 routes for `fetch`, `update`, `insert` and `delete` movies in different ways. Let's deal with them.

```zsh
curl --location 'http://127.0.0.1:4000/v1/movies' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer YCXX5FG6K2ISWOC4SZMRJKAPVA'
```

```json
{
    "metadata": {
        "current_page": 1,
        "page_size": 10,
        "first_page": 1,
        "last_page": 1,
        "total_records": 4
    },
    "movies": [
        {
            "id": 1,
            "title": "avengers",
            "year": 2012,
            "runtime": "2 mins",
            "genres": [
                "sci-fi"
            ],
            "version": 3
        },
        {
            "id": 2,
            "title": "Black Panther",
            "year": 2018,
            "runtime": "134 mins",
            "genres": [
                "action",
                "adventure"
            ],
            "version": 1
        },
        {
            "id": 3,
            "title": "Deadpool",
            "year": 2016,
            "runtime": "108 mins",
            "genres": [
                "action",
                "comedy"
            ],
            "version": 1
        },
        {
            "id": 4,
            "title": "The Breakfast Club",
            "year": 2000,
            "runtime": "96 mins",
            "genres": [
                "drama"
            ],
            "version": 13
        }
    ]
}
```

```zsh
curl --location 'http://127.0.0.1:4000/v1/movies/3' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer YCXX5FG6K2ISWOC4SZMRJKAPVA'
```

```json
{
    "movie": {
        "id": 3,
        "title": "Deadpool",
        "year": 2016,
        "runtime": "108 mins",
        "genres": [
            "action",
            "comedy"
        ],
        "version": 1
    }
}
```

```zsh
curl --location 'http://127.0.0.1:4000/v1/movies' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer YCXX5FG6K2ISWOC4SZMRJKAPVA' \
--data '{
    "title": "The Avengers 2",
    "year": 2012,
    "runtime": "125 mins",
    "genres": [
        "action",
        "comedy",
        "western"
    ]
}'
```

With this request you will get error response:

```json
{
    "error": "your user account doesn't have the necessary permissions to access this resource"
}
```

Need to insert an row in `users_permissions` table for current user to have access to write. After that try again:

```json
{
    "movie": {
        "id": 5,
        "title": "The Avengers 2",
        "year": 2012,
        "runtime": "125 mins",
        "genres": [
            "action",
            "comedy",
            "western"
        ],
        "version": 1
    }
}
```

```zsh
curl --location --request PUT 'http://127.0.0.1:4000/v1/movies/5' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer YCXX5FG6K2ISWOC4SZMRJKAPVA' \
--data '{
    "runtime": "130 mins"
}'
```

```json
{
    "movie": {
        "id": 5,
        "title": "The Avengers 2",
        "year": 2012,
        "runtime": "130 mins",
        "genres": [
            "action",
            "comedy",
            "western"
        ],
        "version": 2
    }
}
```

```zsh
curl --location --request DELETE 'http://127.0.0.1:4000/v1/movies/2' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer YCXX5FG6K2ISWOC4SZMRJKAPVA'
```

```json
{
    "message": "movie deleted successfully"
}
```
