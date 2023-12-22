FROM golang:latest AS builder

WORKDIR /app

COPY /go.mod /go.sum ./
RUN go mod download

COPY / ./

RUN CGO_ENABLED=0 GOOS=linux go build -o pub-server .

FROM alpine
COPY --from=builder /app/pub-server /pub-server
COPY --from=builder /app/ord.json /ord.json

EXPOSE 8081

ENTRYPOINT ["/pub-server"]