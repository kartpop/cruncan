## Running Grafana stack

When you want to test OTEL integration, you can run the Grafana stack with the following command:

```shell
docker compose up -d
```

You can now access the dashboard at http://localhost:3000. The default credentials are admin/admin.

Under explore you can view the metrics, traces and logs when you run an instance of the application with
the following environment variables:

```shell
export OTEL_EXPORTER_OTLP_INSECURE=true
export OTEL_SERVICE_NAME=your-app-name
export OTEL_EXPORTER_OTLP_ENDPOINT=http://127.0.0.1:4318
export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
```

### Viewing logs

In Grafana select Logs/Loki and then select Label: job and select value: your-app-name. You can now view the logs if you run the query.

### Viewing traces

In Grafana select Traces/Tempo. In traceQl paste:

```text
{resource.service.name="your-app-name"}
```

When you click Run Query you can now view the traces.


### Viewing metrics

In Grafana select Explore and then select the data source: prometheus.

Select a metric. You can now view the metrics if you run the query.