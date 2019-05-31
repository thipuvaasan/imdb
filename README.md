## IMDB search catalogue backend

### Endpoints

1. POST `/v1/add/user`

There are two roles, namely: `admin` and `user`.
`Role` is optional in the request body. If unspecified, default role is `user`. If `admin` role is specified, the request maker must be an admin.

Example request body:

```
{
    "name": "Prince Raj",
    "email": "pnc.raj@gmail.com",
    "user_name": "pnc_raj",
    "user_password": "alpha_Imdb"
    "role": "admin"
}
```

Example response on success:
status code: 201
body:

```
{
    "message": "user created successfully"
}
```
2. DELETE `/v1/remove/user`

This endpoint deletes an existing user. Currenlty, admins can delete any profile and user can delete only their profile.
This request body expects an emailID. The user with that emailID will get deleted.

Example request body:

```
{
    "email": "pnc.raj@gmail.com"
}
```

Example response:
status code: 200
body:

```
{
    "message": "user deleted successfully"
}
```

3. POST `/v1/add/movie`

This endpoint adds a new movie in the movie database. Only admins can add a movie. `name`, `99popularity`, `director`, `genre` are required fields.

Example request:
status code: 201

body:

```
{
    "name": "Once upon a time",
    "99popularity": 83,
    "director": "Ajay Devgan",
    "genre": [
        "Comedy",
        "Music",
        "Action"
    ],
    "imdb_score": 8.3
}
```

Example response:

```
{
    "message": "movie added successfully"
}
```

4. DELETE `/v1/remove/movie`

This endpoint deletes a movie from the database. The `movie_id` must be present as a URL param in the request.

Example request:

`DELETE: http://localhost:8000//v1/remove/movie?movie_id=AWsI0f0KI22c2BCr6GxK`

Example response: 
status code: 200
body:

```
"message": "movie deleted successfully"
```

5. PUT `/v1/update/movie`

This endpoint updates a given movie record. `name`, `99popularity`, `director`, `genre` and `movie_id` are required fields.

Example request:

```
{
    "movie_id": AWsI0f0KI22c2BCr6GxK
    "name": "Avengers: Endgame",
    "99popularity": 99,
    "director": "I dont know",
    "genre": [
        "Action"
    ],
    "imdb_score": 9.6
}
```

Example response: 

Example response: 
status code: 200
body:

```
{
    "message": "movie updated successfully"
}
```
6. GET `/v1/get/movie`

This endpoint can be used to search the movies. Accepted URL params are:
    a. `name`
    b. `director`
    c. `genre`
    d. `99popularity`
    e. `imdb_score`
The endpoint also supports pagination. `from` and `size` can be used for pagination. The default value for `from` is 0 and `size` is 20. I have put a cap of 100 on `size`.

Example request:

`GET: http://localhost:8000//v1/get/movie?name=star&size=3`

Example response:
status code: 200

```
{
    "message": "request successful",
    "movies": [
        {
            "movie_id": "AWsH4qrxuDNiuUUjhaC6",
            "name": "Star Trek : The Next Generation",
            "99popularity": 88,
            "director": "Cliff Bole",
            "genre": [
                "Action",
                "Adventure",
                "Sci-Fi"
            ],
            "imdb_score": 8.8
        },
        {
            "movie_id": "AWsH4qrxuDNiuUUjhZ_c",
            "name": "Star Wars",
            "99popularity": 88,
            "director": "George Lucas",
            "genre": [
                "Action",
                "Adventure",
                "Fantasy",
                "Sci-Fi"
            ],
            "imdb_score": 8.8
        },
        {
            "movie_id": "AWsH4qrxuDNiuUUjhZ_h",
            "name": "Star Trek",
            "99popularity": 86,
            "director": "Marc Daniels",
            "genre": [
                "Adventure",
                "Sci-Fi"
            ],
            "imdb_score": 8.6
        }
    ]
}
```
