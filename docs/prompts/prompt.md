Help me write a detailed clear markdown documentation for the repo attached. 


Instructions:
- The project attempts to mimic the architecture that I worked on in a private repo. In order to preserve the learnings during that role as a senior software engineer, this project was created to encapsulate all my learnings related to developing scalable microservices with best practices. Therefore the naming of the services is intentionally dummy - one, two etc. The documentation should go into as much detail as possible related to:
     - code structure
     - golang best practices
     - injecting secretes via config yaml, .env etc.
     - clean application setup via main.go - setting up all configs and dependencies here and injecting them appropriately into various structs
     - logging best practices
     - tracing best practices - more generally, observability using otel library and grafana stack - prometheus, grafana etc.
     - details of the otel setup, metrics, observability etc. - span, tracer, etc.
     - twitter snowflake to generate distributed unique ids across many replicas or services
     - interservices communication:
            - kafka client, kafka consumer, franz-go framework
            - grpc etc. (mock setup in /reference folder)
     - designing REST APIs, graphql apis, consuming third party services
     - authentication, authorization using accesstokens, auth id generator
     - database setup, database migrations, using gorm ORM etc.
     - integration tests, unit tests, module level tests using cucumber godog framework
     - using gherkin fundamentals in testing - given, when, then framework and setting up those tests
     - mocking kafka, databases, third party apis (wiremock) in tests etc.
     - mocking using stubs etc.
     - containerizing applications using Docker
     - distributed retry mechanism using redis redsync
     - setting up cronjobs in AWS to run database cleanup tasks periodically

- document should be as detailed as possible; a new developer should quickly be able to know the workflows, should be able to navigate the codebase easily
- having worked on this project a while ago, this documentation should also act like a refresher for me to know all the important pieces and how they work
- explain all the workflows clearly
- for each workflow, explain how the chain of execution passes through different modules in the code
- IMPORTANT - Create an architecture and workflow diagram using mermaid for the microservices architecture - how backend services one, two commuincate with each other and with external systems - kafka, rest api, graphql, grpc etc.