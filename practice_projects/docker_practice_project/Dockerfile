FROM golang:1.24rc2-bullseye

WORKDIR /app

COPY . .
RUN go mod download

COPY go.mod go.sum ./

RUN GOOS=linux go build -o /main

EXPOSE 8080

CMD ["/main"]