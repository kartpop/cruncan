FROM golang:1.22.2-alpine3.19 AS builder

ENV GH_TOKEN=<github-token-here>

WORKDIR /app

RUN apk add --no-cache ca-certificates git gcc musl-dev tzdata librdkafka-dev

ARG GH_TOKEN
RUN git config --global url."https://${GH_TOKEN}:x-oauth-basic@github.com/kartpop".insteadOf "https://github.com/kartpop"

COPY . .
RUN go mod download
RUN go mod tidy
RUN go build -tags timetzdata,musl,dynamic -o ./out/consumer ./cmd/consumer


FROM alpine:3.19.1 AS consumer
RUN apk add --no-cache ca-certificates librdkafka
COPY --from=builder /app/out/consumer /app/main
COPY --from=builder /app/config /config
ENTRYPOINT ["/app/main"]
EXPOSE 8082
CMD ["-http", ":8082"]
