- Crud application (refer updated Subtasks section below for latest)
  1. POST request: rest api handler + unique auth id generator + save to PostgresDB + publish message to kafka + godog/cucumber component test
       - add open telemetry logging/monitoring/metrics for this request using Prometheus/Grafana stack
  2. Kafka consumer: consumer handler + save to PostgresDB + send grpc request to payments service + make client request to 3rd party API + godog/cucumber component test
       - Dockerfile
  3. Retry job: retry failed client requests in step 2 + use redsync for locking so that multiple replicas don't send same client request + godog/cucumber component test
  4. Cron job: cleanup old jobs (client request is completed and successful in step2) + godog/cucumber component test
       - main variation: Kubernetes cron job
  5. GET request: graphql request + authentication middleware + fetch from DB + godog/cucumber component test
  6. Integrate with example payment gateway: send client request + maintain state in db
  7. Kubernetes + Helm: for all services



Subtasks

1. POST request
     - go mod
     - cmd/first/main.go
     - config/config.go & config.yaml
     - auth id generator
     - db/repository & db/migration/flyway
     - Dockerfile & docker-compose
     - kafka message producer
     - cucumber/godog tests feature
     - tests init and steps
     - prometheus/grafana stack docker-compose
     - logging/traces/metrics

2. Kafka consumer + grpc send request
     - kafka consumer
     - grpc send request
     - save to db

3. grpc receive request + send http client request

3b.  graphql
     - 1. graphql query to fetch data....related to http query created in 1. POST
     - 1. graphql mutation to register user preferences
     

4. reverse flow
     - 3. get http client response + send grpc response
     - 2. receive grpc response + send kafka message (part of same kafka handler as in 2.)
     - 1. kafka consumer for second topic + save to db


Priority
- write UT for kafkahandler in two
- access token validate and verify
