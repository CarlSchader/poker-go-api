FROM golang:alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY application/ ./application/
COPY services/server/ ./server

WORKDIR /usr/src/app/server

RUN go build -v -o /usr/local/bin/app ./...

FROM golang:alpine

WORKDIR /

COPY --from=builder /usr/local/bin/app /usr/local/bin/app

ENV PORT=8080
ENV MONGODB_URI=mongodb://mongo:27017
ENV DB_NAME=poker
ENV RANKS_COLLECTION_NAME=ranks
ENV CACHE_COLLECTION_NAME=cache
ENV POCKETS_COLLECTION_NAME=pockets
ENV CALCULATION_TIMEOUT=120

ENTRYPOINT [ "app" ]
