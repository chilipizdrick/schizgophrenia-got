FROM ubuntu:latest

WORKDIR /app

COPY . /app

RUN go mod tidy && go mod vendor
RUN go mod verify
RUN go build -o ./bin/app

CMD ["./bin/app"]