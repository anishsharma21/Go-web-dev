# Docker, Managed Databases, and Cloud Hosting

The purpose of this project is to learn a lot. For starters, I want to learn how to manage a Postgres database, either with Docker, or more likely with a Managed Database Provider (MDS) like Supabase or Cloud SQL. I also want to learn the differences between cloud providers, like Google Cloud, DigitalOcean, Railway, Fly.io, AWS, and more. Side goals to this include learning how to setup CI/CD so I can establish a full development cycle. I want to also setup a development environment, potentially with docker compose and a locally run postgres db instance running in a volume.

I'd like to learn about logging, analytics, payments, and authentication, but that might spill over into another project.

## Setup steps

I started with a simple go backend server. I created the Dockerfile for the go backend. Then I setup the docker-compose file to run the postgres database in a persistent db volume. To connect to the database in the go backend, I needed to use the following command to get the right postgres driver:

```bash
go get -u github.com/lib/pq
```

To run the postgres database with the backend code locally, navigate to the root of your directory in your terminal, and run `docker-compose up --build`, with the optional `-d` flag if you want to run in a detached state where the containers run in the background rather than in the terminal where you can view the logs. You can then run `docker-compose down` to stop both the server and database container, but the data will be persisted unless you remove the volume for the database itself.

## Run tests

To run tests locally, ensure that you have run the docker-compose command so the tests can connect to the local database. Tests can be run using `go test ./...` command with the option `-v` flag for verbose output of the test outputs. The `./...` is important as it will recursively look for all test files in directories and subdirectories. You can also specify the path directly. It's also a good idea to run `go clean -cache` since Go will cache the results of tests if code doesn't change.

## Local Development

Run `docker compose up -d` to start the local database, and then run `air`. If you don't have `air` installed locally, run `go install github.com/air-verse/air@latest`.
