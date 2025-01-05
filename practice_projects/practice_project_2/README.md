# Go Web Dev - Practice Project 2

## Second attempt at learning the Go + HTMX stack (with other relevant technologies)

### Commands for the docker container that I will probably need

For sqlite driver:

```bash
go get github.com/mattn/go-sqlite3
```

For templ:

```bash
go get github.com/a-h/templ
```

For tailwind executable:

```bash
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
chmod +x tailwindcss-macos-arm64
mv tailwindcss-macos-arm64 tailwindcss
```

To start tailwind builder and watcher:

```bash
./tailwindcss -i ./public/css/input.css -o ./public/css/output.css -w
```

Tailwindcss build for production:

```bash
./tailwindcss -i ./public/css/input.css -o ./public/css/output.css -m
```

Might need to also specify the GOPATH... not sure though
