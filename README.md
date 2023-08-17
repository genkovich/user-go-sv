## Project deployment:

* Copy from env.dist to .env
* Fill in the .env file with the correct data, for example:

```dotenv
PSQL_DSN="postgres://postgres:postgres@localhost/postgres?sslmode=disable"
```

* Run migrations:

```bash
make migration-up
```

* Run the command to start the project:

```bash
go run ./cmd
```