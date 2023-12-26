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
