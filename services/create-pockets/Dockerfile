FROM golang:alpine AS populator

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY application/ ./application/
COPY services/create-pockets/populate/ ./populate

WORKDIR /usr/src/app/populate

RUN go build -v -o /usr/local/bin/app ./...

FROM golang:alpine

WORKDIR /

COPY --from=populator /usr/local/bin/app /usr/local/bin/app
COPY pockets.json /pockets.json

ENV MONGODB_URI=mongodb://mongo:27017
ENV DB_NAME=poker
ENV POCKETS_COLLECTION_NAME=pockets

ENTRYPOINT [ "app", "/pockets.json" ]