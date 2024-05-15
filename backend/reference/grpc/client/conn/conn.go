package conn

import (
	"context"
	"time"

	// grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	// grpc_timeout "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GenericGRPCConn struct {
	Conn *grpc.ClientConn
}

type GRPCRetryConfig struct {
	RetryCount             int `mapstructure:"RETRY_COUNT"`
	PerRetryTimeoutSeconds int `mapstructure:"PER_RETRY_TIMEOUT_SECONDS"`
	InitialBackoffMillis   int `mapstructure:"INITIAL_BACKOFF_MILLIS"`
}

func NewGenericGRPCConnWithRetry(retryConfig GRPCRetryConfig, endpoint string, t time.Duration) *GenericGRPCConn {
	// opts := []grpc_retry.CallOption{
	// 	grpc_retry.WithBackoff(grpc_retry.BackoffExponential(time.Duration(retryConfig.InitialBackoffMillis) * time.Millisecond)),
	// 	grpc_retry.WithMax(uint(retryConfig.RetryCount)),
	// 	grpc_retry.WithPerRetryTimeout(time.Duration(retryConfig.PerRetryTimeoutSeconds) * time.Second),
	// }

	// timeout := t
	// if t <= 0 {
	// 	timeout = 10 * time.Second
	// }

	conn, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// TODO: Replace newrelic grpc interceptor with opentelemetry
		//grpc.WithChainUnaryInterceptor(nrgrpc.UnaryClientInterceptor, grpc_timeout.UnaryClientInterceptor(timeout), grpc_retry.UnaryClientInterceptor(opts...)),
	)

	if err != nil {
		panic(err)
	}

	return &GenericGRPCConn{
		Conn: conn,
	}
}

func NewGenericGRPCConnWithContext(ctx context.Context, endpoint string, t time.Duration) *GenericGRPCConn {
	// timeout := t
	// if t <= 0 {
	// 	timeout = 10 * time.Second
	// }

	conn, err := grpc.DialContext(
		ctx,
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// TODO: Replace newrelic grpc interceptor with opentelemetry
		//grpc.WithChainUnaryInterceptor(nrgrpc.UnaryClientInterceptor, grpc_timeout.UnaryClientInterceptor(timeout)),
		//grpc.WithStreamInterceptor(nrgrpc.StreamClientInterceptor),
	)

	if err != nil {
		panic(err)
	}

	return &GenericGRPCConn{
		Conn: conn,
	}
}

func (client *GenericGRPCConn) Close(ctx context.Context) error {
	return client.Conn.Close()
}
