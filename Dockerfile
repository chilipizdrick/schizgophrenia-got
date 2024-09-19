FROM golang:1.22-bookworm

WORKDIR /app

COPY . /app

RUN go mod tidy && go mod vendor
RUN go mod verify
RUN go build -ldflags "-s -w" -o ./bin/app

CMD ["./bin/app"]
