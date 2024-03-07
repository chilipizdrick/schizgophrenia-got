FROM ubuntu:latest

WORKDIR /app

COPY . /app

RUN go mod download
RUN go mod verify
RUN go build -o ./bin/app

CMD ["./bin/app"]