package tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
	"github.com/kartpop/cruncan/backend/reference/grpc/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ctxKey struct{}

type ctxValue struct {
	transferred int64
	success     bool
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		TestSuiteInitializer: InitializeSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeSuite(ctx *godog.TestSuiteContext) {
	ctx.ScenarioContext().Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ctx = context.WithValue(ctx, ctxKey{}, &ctxValue{})

		return ctx, nil
	})

	ctx.ScenarioContext().After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		//ctxData := ctx.Value(ctxKey{}).(*ctxValue)

		// cleanup db etc. using ctxData

		return ctx, nil
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a transaction request with (\d+) and (\d+) is sent to the server for transfer from "([^"]*)" to "([^"]*)"$`, aTxnReqSentToServer)
	ctx.Step(`^the server processes and sends a transaction response with "([^"]*)" and (\d+) amount$`, serverProcessesTxnReq)
}

func aTxnReqSentToServer(ctx context.Context, amount int64, interest int64, sourceAccountId, targetAccountId string) (context.Context, error) {
	conn, err := grpc.Dial("localhost:8443", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return ctx, err
	}

	req := &model.TransactionRequest{
		SourceAccountId: sourceAccountId,
		TargetAccountId: targetAccountId,
		Amount:          amount,
		Interest:        interest,
	}

	resp, err := model.NewPaymentsServiceClient(conn).DoTransaction(context.Background(), req)
	if err != nil {
		return ctx, err
	}

	ctxData := ctx.Value(ctxKey{}).(*ctxValue)
	ctxData.transferred = resp.Transferred
	ctxData.success = resp.Success

	return ctx, nil
}

func serverProcessesTxnReq(ctx context.Context, successStr string, transferred int64) (context.Context, error) {
	success, err := strconv.ParseBool(successStr)
	if err != nil {
		return ctx, err
	}

	ctxData := ctx.Value(ctxKey{}).(*ctxValue)
	if ctxData.transferred != transferred || ctxData.success != success {
		return ctx, fmt.Errorf("expected response: %v, %v, got: %v, %v", success, transferred, ctxData.success, ctxData.transferred)
	}

	return ctx, nil
}
