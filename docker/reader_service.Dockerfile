FROM golang:1.19-alpine

WORKDIR /app

ENV CONFIG=docker

COPY .. /app

RUN go get -u github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
RUN go mod tidy
RUN go mod download


ENTRYPOINT CompileDaemon --build="go build -o main reader_service/cmd/main.go" --command=./main