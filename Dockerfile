FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o bin/ ./app/...

FROM alpine AS aiko-release

WORKDIR /

COPY --from=build-stage /app/bin/aiko /app/aiko

EXPOSE 2112

ENTRYPOINT ["/app/aiko"]
