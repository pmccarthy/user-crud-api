# User CRUD API

## Introduction
An API Server supporting basic CRUD commands for user management backed by a PostgreSQL database

### Prerequisites
* go v1.20+

## Running locally
Prior to running the server locally you will first need to specify configuration details of a target PostgreSQL v14.x instance using a local `.env` file

### Create local .env file
```bash
touch .env
```

### Configure local .env file
Example:
``` bash
DB_HOST=localhost
DB_PORT=5455
DB_USER=postgresUser
DB_PASS=postgresPW
DB_NAME=users
DB_SSLMODE=disable
```

### Using a Docker Postgres Instance
Docker is a good option for running PostgreSQL locally assuming you don't already have an instance in place already.

#### Pull down PostgreSQL image
```bash
docker pull postgres
```
#### Run docker PostgreSQL image locally
Note: The below command will work with the default values set in the `.env` file of this project. Should you choose to modify values in the `.env` file please ensure that you also modify them here:
```bash
docker run --name myPostgresDb -p 5455:5432 -e POSTGRES_USER=postgresUser -e POSTGRES_PASSWORD=postgresPW -e POSTGRES_DB=users -d postgres
```

### Run the server
From the root of the project, run the following:
```bash
go run .
```

## Examples

### Create a new user
```bash
curl --location 'http://localhost:8080/api/create_user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "user1",
    "password": "password",
    "email": "email@example.com"
}'
```

### Get list of users
```bash
curl --location 'http://localhost:8080/api/get_users'
```

### Get user by ID
```bash
ID=1; curl --location "http://localhost:8080/api/get_user/+${ID}"
```

### Delete a user
```bash
ID=1; curl --location --request DELETE "http://localhost:8080/api/delete_user/+${ID}"
```
