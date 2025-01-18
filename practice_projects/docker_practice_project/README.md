# Docker Practice Project for Go web dev

The purpose of this practice project is to pair my current experience in Go web development with Docker and deploying my application. The tech stack is almost completely the same as practice project 2, except I am opting for the standard library text/template instead of Templ to reduce dependencies and complexity of build process.

## Setup process

Tailwind executable:

```bash
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
chmod +x tailwindcss-macos-arm64
mv tailwindcss-macos-arm64 tailwindcss
```

SQlite driver:

```bash
go get github.com/mattn/go-sqlite3
```
