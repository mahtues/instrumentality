FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o instrumentality ./cmd/app

EXPOSE 80

CMD ./instrumentality
