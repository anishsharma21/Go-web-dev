FROM golang:bullseye AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN GOOS=linux go build -o /main

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /main /main

EXPOSE 8080

CMD ["/main"]