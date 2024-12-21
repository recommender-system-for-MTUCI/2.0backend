# Specification special for Daria

Server is starting on __localhost:8080__


## API

##### User

* `POST /api/registration` - sign up
* `GET /api/profile` - get profile
* `PATCH /api/update_password` - update password
* `DELETE /api/delete_user` - delete profile
* `POST /api/login` - sign in
* `POST /api/accept_email` - accept email
* `GET /api/refresh` - update jwt token

##### Film

* `GET /api/recommend_system` - return main page with films
* `GET /api/favorites` - return favorites film
* `GET /api/movie/:id` - return film by ID
* `DELETE /api/favorites/:id` - delete film from favorites by ID
* `POST /api/comment/:id` - add new comment by film ID
* `DELETE /api/comment/:id` - delete comment by comment ID
* `POST /api/favorites/:id` - add film to favorites by film ID
* `GET /api/genres` - return all unique genres
* `GET /api/:genre/:page` - return filtration by genres
* `GET/api/film/:name` - return film by name
* `POST /api/comments/:id` - return comments by film ID

### Sign up 

Handler : `POST /api/registration`.

Registration is done by email and password. 

Every email should be unique.

Password should be more than 8 characters.

Request form 

```
POST /api/registration
Content-Type: application/json
...

{
  "login": "<login>",
  "password": "<password>"
}
```

Possible request status:

- `200` - user successful created;

Response form:

```
200 OK
Content-Type: application/json
...

{
      "access_token": "<access_token>",
      "refresh_token": "<refresh_token>"
    }
```

- `400` - wrong request form;
- `500` - server error.

### Get profile

Handler: `GET /api/profile`

Return user profile

Request form:

```
GET /api/profile
Authorization: Bearer<access_token>
```
Possible request status:
- `200` - profile successful returned;

Response form:
```
200 OK
Content-Type: application/json
{
    "login":<login>
}
```
- `500` - server error;
- `403` - user didn`t accept email.

### Update password

Handler : `PATCH /api/update_password`

This request update user password.

Password should be more than 8 characters.

Request form:

```
PATCH /api/update_password
Content-Type: application/json
Authorization: Bearer<access_token>
...

{
    "old_password":"<old_password>"
    "new_password":"<new_password>"
}
```
Possible request status:

- `204` - password was updated successful;
- `500` - server error;
- `403` - user didn`t accept email;
- `400` - wrong request.

### Delete profile

Handler: `DELETE /api/delete_user`

This request delete user profile

Request form:

```
DELETE /api/delete_user
Authorization: Bearer<access_token>
```

Possible request status:

- `200` - user was deleted successful;
- `500` - server error;

### Sign in

Handler: `POST /api/login`

This request for authorization.

Password should be more than 8 characters.

Request form:
```
POST /api/login
Content-Type: application/json
...

{
  "login": "<login>",
  "password": "<password>"
}
```
Possible request status:

- `200` - user successful signed in;
 ```
200 OK
Content-Type: application/json
...

{
      "access_token": "<access_token>",
      "refresh_token": "<refresh_token>"
}
```
- `400` - wrong request form;
- `500` - internal server error;
- `403` - user didn`t accept email.

### Accept email

Handler: `POST /api/accept_email`

This is request for accept email

Request form:

```
POST /api/accept_email
Content-Type: application/json
Authorization: Bearer<access_token>
...

{
    "code":code
}
```

Possible response status:

- `200` - email successful accepted;
- `500` - server error;
- `400` - wrong request form;

### Update jwt token

Handler: `GET /api/refresh`

When access token lifetime is over, you should use this request

If refresh token lifetime is over, you should ask user to sign in again

Request form:

```
GET /api/refresh
Authorization: Bearer<refresh_token>
```

Possible response status:

- `200` successful update tokens; 

Response form:

```
200 OK
Content-Type: application/json
...

{
      "access_token": "<access_token>",
      "refresh_token": "<refresh_token>"
    }
```
- `400` - wrong request form;
- `401` - wrong refresh token;
- `500` - server error.

### Main page (_NEED FIX JSON_)

Handler: `GET /api/recommend_system`

Return main page with 20 popular films

Request form:

```
GET /api/recommend_system
```

Possible response status:

- `200` - returned 20 films;

Response form:

```
200 OK
Content-Type: application/json
...

[
    {
        "film_id":"<film_id>",
        "name":"<name>",
        "rating":<rating>(float)
    },
    {
        "film_id":"<film_id>",
        "name":"<name>",
        "rating":<rating>(float)
    },
]

```

- `500` - server error.

### Favorites

Handler: `GET /api/favorites`

Return users favorites
    
