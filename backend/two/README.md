# two
A service which consumes messages from onerequest kafka topic and sends a POST request on the threerequest API.

## Setup

- ensure kafka and grafana stack are running
- setup dependencies (wiremock etc.)
    ```shell
    docker compose up -d
    ```
- run app
    ```shell
    go run cmd/consumer/main.go
    ```
- run tests
    ```shell
    go test -v ./... -count 1  --shuffle=on --parallel=16
    ```