## Project deployment:

### Download and install:
* Clone project
```bash
git clone git@github.com:genkovich/user-go-sv.git .
```

* Install dependencies:

```bash
go mod download
```
---

### Setup database

* Copy from ./docker/env.dist to ./docker/.env
* Fill in the .env file with the correct data, for example:
```dotenv
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=postgres
POSTGRES_PORT=5432
```

* Setup database:
```bash
make setup-local-db
```
---
### Run project

* Copy from env.dist to .env
* Fill in the .env file with the correct data, for example:

```dotenv
PSQL_DSN="postgres://postgres:postgres@localhost/postgres?sslmode=disable"
JWT_SECRET="TFtGJbQ9drDOSNrfRjJFbiUy55AF4/dXyOj+7lvmYBDPZB8BV97vEVwdxSgy0SDoFQlIOzL1oAh9VVT5oK82gQ=="
```
* Run the command to start the project:

```bash
go run ./cmd
```