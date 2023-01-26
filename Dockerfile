FROM golang:1.19-alpine3.17 AS builder
WORKDIR /go/app
COPY . ./
RUN go mod download
RUN mkdir bin
RUN cd cmd/main/ && go build -o ../../bin/go-redis-mapper

FROM golang:1.19-alpine3.17
WORKDIR /go
COPY --from=builder /go/app/bin/go-redis-mapper ./
ENTRYPOINT ["./go-redis-mapper"]
