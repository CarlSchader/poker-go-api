FROM golang:alpine AS generator

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY application/ ./application/
COPY services/create-ranks/generate/ ./generate

WORKDIR /usr/src/app/generate

RUN go build -v -o /usr/local/bin/app ./...
RUN app /ranks.json

FROM golang:alpine AS populator

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY application/ ./application/
COPY services/create-ranks/populate/ ./populate

WORKDIR /usr/src/app/populate

RUN go build -v -o /usr/local/bin/app ./...

FROM golang:alpine

WORKDIR /

COPY --from=populator /usr/local/bin/app /usr/local/bin/app
COPY --from=generator /ranks.json /ranks.json

ENV MONGODB_URI=mongodb://mongo:27017
ENV DB_NAME=poker
ENV RANKS_COLLECTION_NAME=ranks
ENV CACHE_COLLECTION_NAME=cache
ENV BATCHES=20

ENTRYPOINT [ "app", "/ranks.json" ]