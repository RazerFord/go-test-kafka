ARG GO_VERSION="1.22"
ARG APP_TYPE="client"

######################################################

FROM golang:${GO_VERSION} AS app-client

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/client/...

######################################################

FROM golang:${GO_VERSION} AS app-server

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/server/...

######################################################

FROM app-${APP_TYPE} AS app
CMD ["/usr/local/bin/app"]
