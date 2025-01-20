# Docker Practice Project for Go web dev

The purpose of this practice project is to pair my current experience in Go web development with Docker and deploying my application. The tech stack is almost completely the same as practice project 2, except I am opting for the standard library text/template instead of Templ to reduce dependencies and complexity of build process.

## Setup process

### Database (`goose` and sqlite)

Install SQlite driver

```bash
go get github.com/mattn/go-sqlite3
```

Github link for `goose` at this link [here](https://github.com/pressly/goose?tab=readme-ov-file)
Documentation for `goose` at this link [here](https://pressly.github.io/goose/)

Get `goose` for database handling and migrations

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Create first table

```bash
goose create add_books_table sql
```

Variables need to be set, either as exported variables or in a `.env` file. I've gone with exported variables for this project.

```bash
export GOOSE_DRIVER=sqlite3
export GOOSE_DBSTRING=./app.db
export GOOSE_MIGRATION_DIR=./migrations
```

After setting these environment variables, you can simple run `goose up` or `goose down` to apply the migration.

### Tailwind

Install the executable

```bash
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
chmod +x tailwindcss-macos-arm64
mv tailwindcss-macos-arm64 tailwindcss
```

Initialise `tailwind.config.js` file with `./tailwindcss init`. Then, created an `input.css` file in the `static/css` directory. Then run command to build `output.css` file.

```bash
./tailwindcss -i static/css/input.css -o public/css/output.css
```

This output file is not tracked and needs to be rebuilt on deployment. For production, this file should be rebuilt with the `--minify` flag.

### Docker commands

First build the Docker image (cd into the root of this project)

```bash
docker build -t bookapp .
```

Then run the container, making sure to map your local port to 8080 if you want to access the site locally

```bash
docker run -p 8080:8080 bookapp
```
