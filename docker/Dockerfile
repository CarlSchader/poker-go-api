FROM golang:1.17.8-bullseye

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY src/ .
RUN go build -v -o /usr/local/bin/app ./...

ENV PORT=8080

EXPOSE 8080

ENTRYPOINT ["app"]