FROM golang:1.17.8-bullseye AS generator

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY services/create-ranks/generate ./
RUN go build -v -o /usr/local/bin/app ./...
RUN app /data/ranks.json

FROM golang:1.17.8-bullseye

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY --from=generator /data/ranks.json ./
COPY services/create-ranks/populate ./

RUN go build -v -o /usr/local/bin/app ./...

ENV MONGODB_URI=localhost:27017
ENV DB_NAME=poker
ENV COLLECTION_NAME=ranks

ENTRYPOINT [ "app", "ranks.json" ]
