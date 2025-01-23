# Stage 1: Build the Go application
FROM golang:bullseye AS builder

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN GOOS=linux go build -o main

# Stage 2: Create a minimal image
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=builder /main /main
COPY --from=builder /internal/templates /internal/templates
COPY --from=builder /static /static

EXPOSE 8080

CMD ["/main"]