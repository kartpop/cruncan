- Crud application
  1. POST request: rest api handler + unique auth id generator + save to PostgresDB + publish message to kafka + godog/cucumber component test
       - add open telemetry logging/monitoring/metrics for this request using Prometheus/Grafana stack
  3. Kafka consumer: consumer handler + save to PostgresDB + make client request to 3rd party API + godog/cucumber component test
  4. Retry job: retry failed client requests in step 2 + use redsync for locking so that multiple replicas don't send same client request + godog/cucumber component test
  5. Cron job: cleanup old jobs (client request is completed and successful in step2) + godog/cucumber component test
       - main variation: Kubernetes cron job
  6. GET request: graphql request + authentication middleware + fetch from DB + godog/cucumber component test
  7. POST payment request: rest api handler + maintain state in PostgresDB + send grpc request to payments service + godog/cucumber component test
       - integrate with example payment gateway: send client request + maintain state in db
  7. Kubernetes + Helm: for all services
