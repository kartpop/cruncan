package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"io"

	"github.com/cucumber/godog"
	"github.com/kartpop/cruncan/backend/one/database/onerequest"
	"github.com/kartpop/cruncan/backend/one/tests/utils"
	kafkaUtil "github.com/kartpop/cruncan/backend/pkg/kafka"
	"github.com/kartpop/cruncan/backend/pkg/model"
	"gorm.io/gorm"
)

type testFixtureKey struct{}

type testFixture struct {
	reqId                string
	status               int
	kafkaClient          *kafkaUtil.Client
	oneRequestKafkaMssgs map[string]model.OneRequest
	gormClient           *gorm.DB
	oneRequestRepo       onerequest.Repository
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		TestSuiteInitializer: InitializeSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeSuite(ctx *godog.TestSuiteContext) {
	var key = testFixtureKey{}

	ctx.ScenarioContext().Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		kafkaClient := utils.InitKafkaClient()
		oneRequestTestConsumer := kafkaClient.NewConsumer(utils.EnvConfig.Kafka.OneRequestTopic.Name)

		gormClient := utils.InitGorm()
		oneRequestRepo := onerequest.NewRepository(gormClient)

		ctx = context.WithValue(ctx, key, &testFixture{
			kafkaClient:          kafkaClient,
			oneRequestKafkaMssgs: make(map[string]model.OneRequest),
			gormClient:           gormClient,
			oneRequestRepo:       oneRequestRepo,
		})

		oneRequestTestConsumer.Start(ctx, &TestKafkaHandler{})

		return ctx, nil
	})

	ctx.ScenarioContext().After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		ctxKeyValue := ctx.Value(key)
		ctxData := ctxKeyValue.(*testFixture)

		ctxData.kafkaClient.Close()

		ctxData.gormClient.WithContext(ctx).Where("1 = 1").Delete(&onerequest.OneRequest{})
		db, er := ctxData.gormClient.DB()
		if er != nil {
			return ctx, er
		}

		return ctx, db.Close()
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a POST request with JSON body from file "([^"]*)" is sent$`, aPostRequestIsSent)
	ctx.Step(`^the request with message from file "([^"]*)" is published to kafka$`, theRequestIsPublishedToKafka)
	ctx.Step(`^the request with message from file "([^"]*)" is saved to database with correct user id$`, theRequestIsSavedToDatabase)
	ctx.Step(`^the response status code is (\d+)$`, responseCodeIs)
}

func aPostRequestIsSent(ctx context.Context, filePath string) (context.Context, error) {
	ctxData := ctx.Value(testFixtureKey{}).(*testFixture)

	reqBody, err := os.ReadFile(filePath)
	if err != nil {
		return ctx, fmt.Errorf("could not read file %s: %v", filePath, err)
	}

	url := "http://" + utils.EnvConfig.Server.Addr + "/one"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return ctx, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ctx, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return ctx, err
	}
	var responseBody map[string]interface{}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return ctx, err
	}
	ctxData.reqId = responseBody["request_id"].(string)
	ctxData.status = resp.StatusCode

	return ctx, nil
}

func theRequestIsPublishedToKafka(ctx context.Context, filePath string) (context.Context, error) {
	reqBody, err := os.ReadFile(filePath)
	if err != nil {
		return ctx, fmt.Errorf("could not read file %s: %v", filePath, err)
	}

	var oneReq model.OneRequest
	if err := json.Unmarshal(reqBody, &oneReq); err != nil {
		return ctx, fmt.Errorf("failed to unmarshal one request, error: %v", err)
	}

	// polling to periodically check if onerequest test kafka handler has processed the message
	timeout := time.After(15 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return ctx, fmt.Errorf("test timed out")
		case <-tick:
			err := verifyMessage(ctx, oneReq)
			if err != nil {
				slog.Default().Error(err.Error())
			}
			if err == nil {
				return ctx, nil
			}
		}
	}
}

func theRequestIsSavedToDatabase(ctx context.Context, filePath string) (context.Context, error) {
	ctxData := ctx.Value(testFixtureKey{}).(*testFixture)

	reqBody, err := os.ReadFile(filePath)
	if err != nil {
		return ctx, fmt.Errorf("could not read file %s: %v", filePath, err)
	}

	var oneReq model.OneRequest
	if err := json.Unmarshal(reqBody, &oneReq); err != nil {
		return ctx, fmt.Errorf("failed to unmarshal one request, error: %v", err)
	}

	oneReqFromDB, err := ctxData.oneRequestRepo.Get(ctx, ctxData.reqId)
	if err != nil {
		return ctx, fmt.Errorf("failed to get one request from database, error: %v", err)
	}

	if oneReqFromDB.UserID != oneReq.UserID {
		return ctx, fmt.Errorf("user id mismatch")
	}

	return ctx, nil
}

func responseCodeIs(ctx context.Context, expectedStatusCode int) (context.Context, error) {
	ctxData := ctx.Value(testFixtureKey{}).(*testFixture)

	if ctxData.status != expectedStatusCode {
		return ctx, fmt.Errorf("expected status code %d, but got %d", expectedStatusCode, ctxData.status)
	}

	return ctx, nil
}

func verifyMessage(ctx context.Context, oneReq model.OneRequest) error {
	ctxData := ctx.Value(testFixtureKey{}).(*testFixture)

	if _, ok := ctxData.oneRequestKafkaMssgs[oneReq.UserID]; !ok {
		return fmt.Errorf("one request not found in kafka messages")
	}

	// verify sample data
	if ctxData.oneRequestKafkaMssgs[oneReq.UserID].UserID != oneReq.UserID {
		return fmt.Errorf("user id mismatch")
	}
	if ctxData.oneRequestKafkaMssgs[oneReq.UserID].Prompt != oneReq.Prompt {
		return fmt.Errorf("prompt mismatch")
	}

	return nil
}

type TestKafkaHandler struct {
}

func (t *TestKafkaHandler) Handle(ctx context.Context, msg []byte, topic string) error {
	ctxData := ctx.Value(testFixtureKey{}).(*testFixture)

	var oneReq model.OneRequest
	if err := json.Unmarshal(msg, &oneReq); err != nil {
		slog.Default().Error(fmt.Sprintf("failed to unmarshal one request, error: %v", err))
		return err
	}

	ctxData.oneRequestKafkaMssgs[oneReq.UserID] = oneReq

	return nil
}
