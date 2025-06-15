# Chirpy


## Eng

### Description

Chirpy is a twitter-like is a back-end server Go program. This project
was designed to learn and practice http-server concepts and notions.
It handles http requests and returns responses both in JSON format.
There is no frontend, it's just the backend logic.
The endpoints and functionality are explained in more detail below.
The program uses PostgreSQL for db management, sqlc to generate Go code
from SQL and pressly/Goose to make database migrations.

### Requirements

To run this program, make sure you have the latest versions of Go
and PostgreSQL installed on your system.

### Set-Up

- Clone the repo to your system and install:

```bash
git clone https://github.com/iegpeppino/chirpy
cd path/to/chirpy
go install ...
```

- Create a PostgreSQL database named chirpy 

- Store the DB URL, a secret key, and the Polka ApiKey (it's made up) in environment variables.

- Make the necessary database migrations:

```bash
cd sql/schema
goose postgres <connection_url> up
```

- Return to root directory, compile and run the code:

```bash
cd ../..
go build -o out && ./out
```

### Functionality

You can create, login and update users using an email and password.
User's password gets hashed and they receive an access token. This token
is kept being refreshed and can be revoked.
There's a webhook to upgrade a user status to 'chirpy_red' through a fictional
api called Polka. 
Loged users can create Chirps (twits) and delete them (if they own them).
You can list all chirps for a user using their user_id as query parameter, or
read a specific chirp using its id.
Most request perform user JWT token validation before processing the response.

### Usage

#### Available Endpoints 

* Create user
    - Method: "POST"
    - Endpoint: /api/users
    - Params: 
        ```json
        {
            email: "user@email.com"
            password: "mypasswordsnot123"
        }
        ```

* Login user
    - Method: "POST"
    - Endpoint: /api/login
    - Params: 
        ```json
        {
            email: "user@email.com"
            password: "mypasswordsnot123"
        }
        ```

* Update User's email and password
    - Method: "PUT"
    - Endpoint: /api/users
    - Params:
        ```json
        {
            email: "user@email.com"
            password: "mypasswordsnot123"
        }
        ```

* Delete all users from db
    * Disclaimer: Your client's platform must be set as "dev"
    - Method: "POST"
    - Endpoint: /admin/reset

* Create chirp (twitt)
    - Method: "POST"
    - Endpoint: /api/chirps
    - Params:
        ```json
        {
            body: "hello there!"
        }
        ```

* List all chirps
    Note: By default it lists all chirps from db, but you can
    pass "user_id" and "order" (asc or desc) as url query parameters
    to list all chirps from a selected user and sort them by date.
    - Method: "GET"
    - Endpoint: /api/chirps

* Get chirp by ID
    - Method: "GET"
    - Endpoint /api/chiprs/{chirpID}

* Delete chirp by ID
    - METHOD: "DELETE"
    - Endpoint: /api/chirps/{chirpID}

* Upgrade user to "chirpy_red_status"
    * Disclaimer: event param must be "user.upgraded" for success
    - Method: "POST"
    - Endpoint: /api/polka/webhooks
    - Params:
        ```json
        {
            event: "some.event"
            data: {
                user_id: "asfj0123jf1093jf"
            }
        }
        ```

* Refresh Access Token
    - Method: "POST"
    - Endpoint: /api/refresh
    - Params:
        ```json
        {
            token: "userJWTtokenstring"
        }
        ```

* Revoke Refresh Access token from User
    - Method: "POST"
    - Endpoint: /api/revoke

* Other less critical endpoints:
    - "GET /api/healthz" : Returns code 200 to signal the server is running
    - "GET /admin/metrics" : Returns number of server hits since it started running