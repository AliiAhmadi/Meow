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

## Health check

Another route implemented is a "healthcheck" to fetch some information about the API (no authentication is required):

```zsh
curl --location 'http://127.0.0.1:4000/v1/healthcheck' \
--header 'Content-Type: application/json'
```

```json
{
    "info": {
        "environment": "development",
        "status": "available",
        "version": "1.0.0"
    }
}
```

## Debug

And the last route of this simple API is `debug`. Use it to see system resource usage and many more data:

```zsh
curl --location 'http://127.0.0.1:4000/debug' \
--header 'Content-Type: application/json'
```

```json
{
    "cmdline": [
        "/tmp/go-build1112461546/b001/exe/api",
        "-smtp-username=4a3738588400bf",
        "-smtp-password=c3534223e2f462"
    ],
    "database": {
        "MaxOpenConnections": 25,
        "OpenConnections": 0,
        "InUse": 0,
        "Idle": 0,
        "WaitCount": 0,
        "WaitDuration": 0,
        "MaxIdleClosed": 0,
        "MaxIdleTimeClosed": 1,
        "MaxLifetimeClosed": 0
    },
    "goroutines": 8,
    "inflight_requests": 0,
    "memstats": {
        "Alloc": 932584,
        "TotalAlloc": 932584,
        "Sys": 8494096,
        "Lookups": 0,
        "Mallocs": 5377,
        "Frees": 401,
        "HeapAlloc": 932584,
        "HeapSys": 3571712,
        "HeapIdle": 1335296,
        "HeapInuse": 2236416,
        "HeapReleased": 1204224,
        "HeapObjects": 4976,
        "StackInuse": 622592,
        "StackSys": 622592,
        "MSpanInuse": 57792,
        "MSpanSys": 65184,
        "MCacheInuse": 14400,
        "MCacheSys": 15600,
        "BuckHashSys": 4005,
        "GCSys": 2915544,
        "OtherSys": 1299459,
        "NextGC": 4194304,
        "LastGC": 0,
        "PauseTotalNs": 0,
        "PauseNs": [
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0
        ],
        "PauseEnd": [
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0
        ],
        "NumGC": 0,
        "NumForcedGC": 0,
        "GCCPUFraction": 0,
        "EnableGC": true,
        "DebugGC": false,
        "BySize": [
            {
                "Size": 0,
                "Mallocs": 0,
                "Frees": 0
            },
            {
                "Size": 8,
                "Mallocs": 134,
                "Frees": 0
            },
            {
                "Size": 16,
                "Mallocs": 1544,
                "Frees": 0
            },
            {
                "Size": 24,
                "Mallocs": 555,
                "Frees": 0
            },
            {
                "Size": 32,
                "Mallocs": 363,
                "Frees": 0
            },
            {
                "Size": 48,
                "Mallocs": 613,
                "Frees": 0
            },
            {
                "Size": 64,
                "Mallocs": 224,
                "Frees": 0
            },
            {
                "Size": 80,
                "Mallocs": 204,
                "Frees": 0
            },
            {
                "Size": 96,
                "Mallocs": 173,
                "Frees": 0
            },
            {
                "Size": 112,
                "Mallocs": 408,
                "Frees": 0
            },
            {
                "Size": 128,
                "Mallocs": 74,
                "Frees": 0
            },
            {
                "Size": 144,
                "Mallocs": 126,
                "Frees": 0
            },
            {
                "Size": 160,
                "Mallocs": 89,
                "Frees": 0
            },
            {
                "Size": 176,
                "Mallocs": 21,
                "Frees": 0
            },
            {
                "Size": 192,
                "Mallocs": 5,
                "Frees": 0
            },
            {
                "Size": 208,
                "Mallocs": 61,
                "Frees": 0
            },
            {
                "Size": 224,
                "Mallocs": 20,
                "Frees": 0
            },
            {
                "Size": 240,
                "Mallocs": 2,
                "Frees": 0
            },
            {
                "Size": 256,
                "Mallocs": 55,
                "Frees": 0
            },
            {
                "Size": 288,
                "Mallocs": 39,
                "Frees": 0
            },
            {
                "Size": 320,
                "Mallocs": 11,
                "Frees": 0
            },
            {
                "Size": 352,
                "Mallocs": 54,
                "Frees": 0
            },
            {
                "Size": 384,
                "Mallocs": 1,
                "Frees": 0
            },
            {
                "Size": 416,
                "Mallocs": 66,
                "Frees": 0
            },
            {
                "Size": 448,
                "Mallocs": 1,
                "Frees": 0
            },
            {
                "Size": 480,
                "Mallocs": 0,
                "Frees": 0
            },
            {
                "Size": 512,
                "Mallocs": 7,
                "Frees": 0
            },
            {
                "Size": 576,
                "Mallocs": 12,
                "Frees": 0
            },
            {
                "Size": 640,
                "Mallocs": 4,
                "Frees": 0
            },
            {
                "Size": 704,
                "Mallocs": 7,
                "Frees": 0
            },
            {
                "Size": 768,
                "Mallocs": 2,
                "Frees": 0
            },
            {
                "Size": 896,
                "Mallocs": 5,
                "Frees": 0
            },
            {
                "Size": 1024,
                "Mallocs": 10,
                "Frees": 0
            },
            {
                "Size": 1152,
                "Mallocs": 13,
                "Frees": 0
            },
            {
                "Size": 1280,
                "Mallocs": 6,
                "Frees": 0
            },
            {
                "Size": 1408,
                "Mallocs": 3,
                "Frees": 0
            },
            {
                "Size": 1536,
                "Mallocs": 12,
                "Frees": 0
            },
            {
                "Size": 1792,
                "Mallocs": 2,
                "Frees": 0
            },
            {
                "Size": 2048,
                "Mallocs": 4,
                "Frees": 0
            },
            {
                "Size": 2304,
                "Mallocs": 3,
                "Frees": 0
            },
            {
                "Size": 2688,
                "Mallocs": 1,
                "Frees": 0
            },
            {
                "Size": 3072,
                "Mallocs": 0,
                "Frees": 0
            },
            {
                "Size": 3200,
                "Mallocs": 1,
                "Frees": 0
            },
            {
                "Size": 3456,
                "Mallocs": 0,
                "Frees": 0
            },
            {
                "Size": 4096,
                "Mallocs": 10,
                "Frees": 0
            },
            {
                "Size": 4864,
                "Mallocs": 2,
                "Frees": 0
            },
            {
                "Size": 5376,
                "Mallocs": 1,
                "Frees": 0
            },
            {
                "Size": 6144,
                "Mallocs": 2,
                "Frees": 0
            },
            {
                "Size": 6528,
                "Mallocs": 2,
                "Frees": 0
            },
            {
                "Size": 6784,
                "Mallocs": 0,
                "Frees": 0
            },
            {
                "Size": 6912,
                "Mallocs": 0,
                "Frees": 0
            },
            {
                "Size": 8192,
                "Mallocs": 4,
                "Frees": 0
            },
            {
                "Size": 9472,
                "Mallocs": 13,
                "Frees": 0
            },
            {
                "Size": 9728,
                "Mallocs": 0,
                "Frees": 0
            },
            {
                "Size": 10240,
                "Mallocs": 0,
                "Frees": 0
            },
            {
                "Size": 10880,
                "Mallocs": 1,
                "Frees": 0
            },
            {
                "Size": 12288,
                "Mallocs": 0,
                "Frees": 0
            },
            {
                "Size": 13568,
                "Mallocs": 0,
                "Frees": 0
            },
            {
                "Size": 14336,
                "Mallocs": 0,
                "Frees": 0
            },
            {
                "Size": 16384,
                "Mallocs": 0,
                "Frees": 0
            },
            {
                "Size": 18432,
                "Mallocs": 1,
                "Frees": 0
            }
        ]
    },
    "timestamp": 1703837147,
    "total_processing_time_microseconds": 591934,
    "total_requests_received": 17,
    "total_responses_send": 16,
    "total_responses_sent_by_status": {
        "200": 7,
        "201": 2,
        "401": 5,
        "403": 1,
        "422": 1
    },
    "version": "1.0.0"
}
```
