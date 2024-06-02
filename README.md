# Project adinata

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Migration
[create migration](https://github.com/golang-migrate/migrate) installed

export POSTGRESQL_URL
```bash
export POSTGRESQL_URL='postgres://username:password@localhost:5432/db_name?sslmode=disable'
```
run migration command
```bash
migrate -database ${POSTGRESQL_URL} -path files/db/migrations up
```




## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```
