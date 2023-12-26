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
