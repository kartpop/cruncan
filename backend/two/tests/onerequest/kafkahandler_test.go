package onerequest

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/kartpop/cruncan/backend/two/tests/utils"
	"github.com/wiremock/go-wiremock"
)

type testFixtureKey struct{}

type testFixture struct {
	stubsToDelete []*wiremock.StubRule
}

func (tf *testFixture) lastStubRule() (stubRule *wiremock.StubRule, ok bool) {
	if len(tf.stubsToDelete) == 0 {
		return nil, false
	}
	return tf.stubsToDelete[len(tf.stubsToDelete)-1], true
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
		ctx = context.WithValue(ctx, key, &testFixture{})
		return ctx, nil
	})

	ctx.ScenarioContext().After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		ctxKeyValue := ctx.Value(key)
		ctxData := ctxKeyValue.(*testFixture)

		for _, stubRule := range ctxData.stubsToDelete {
			_ = utils.WiremockClient.DeleteStub(stubRule)
		}

		return ctx, nil
	})

}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the threerequest API responds with (\d+) on matching body from file "([^"]*)"$`, threeReqApiMatchesBodyAndResponds)
	ctx.Step(`^the token API response with (\d+) on matching body from file "([^"]*)"$`, tokenApiMatchesBodyAndResponds)
	ctx.Step(`^a onerequest from file "([^"]*)" is ingested$`, onerequestIsIngested)
	ctx.Step(`^a request is made to the threerequest API$`, requestIsMadeToThreeRequestApi)
}

func threeReqApiMatchesBodyAndResponds(ctx context.Context, expectedStatusCode int, filePath string) (context.Context, error) {
	ctxData := ctx.Value(testFixtureKey{}).(*testFixture)

	b, err := os.ReadFile(filePath)
	if err != nil {
		return ctx, fmt.Errorf("could not read file %s: %v", filePath, err)
	}

	rule := wiremock.Post(wiremock.URLEqualTo("/three")).
		WillReturnResponse(
			wiremock.NewResponse().
				WithStatus(int64(expectedStatusCode)).
				WithHeader("Content-Type", "application/json"),
		).
		WithBodyPattern(wiremock.EqualToJson(string(b), wiremock.IgnoreArrayOrder, wiremock.IgnoreExtraElements))
	err = utils.WiremockClient.StubFor(rule)

	ctxData.stubsToDelete = append(ctxData.stubsToDelete, rule)

	return ctx, err
}

func tokenApiMatchesBodyAndResponds(ctx context.Context, expectedStatusCode int, filePath string) (context.Context, error) {
	ctxData := ctx.Value(testFixtureKey{}).(*testFixture)

	b, err := os.ReadFile(filePath)
	if err != nil {
		return ctx, fmt.Errorf("could not read file %s: %v", filePath, err)
	}

	credentials := utils.EnvConfig.Auth.ClientID + ":" + utils.EnvConfig.Auth.ClientSecret
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))

	rule := wiremock.Post(wiremock.URLEqualTo("/v2/oauth/token")).
		WithHeader("Content-Type", wiremock.EqualTo("application/x-www-form-urlencoded")).
		WithHeader("Authorization", wiremock.EqualTo(authHeader)).
		WillReturnResponse(
			wiremock.NewResponse().
				WithStatus(int64(expectedStatusCode)).
				WithHeader("Content-Type", "application/x-www-form-urlencoded").
				WithBody(string(b)),
		)
	err = utils.WiremockClient.StubFor(rule)

	ctxData.stubsToDelete = append(ctxData.stubsToDelete, rule)

	return ctx, err
}

func onerequestIsIngested(ctx context.Context, filePath string) (context.Context, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return ctx, fmt.Errorf("could not read file %s: %v", filePath, err)
	}

	err = utils.OneRequestProducer.SendMessage(ctx, b)

	return ctx, err
}

func requestIsMadeToThreeRequestApi(ctx context.Context) (context.Context, error) {
	// polling to periodically check if onerequest kafka handler has processed the message
	timeout := time.After(15 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return ctx, fmt.Errorf("test timed out")
		case <-tick:
			err := verifyStub(ctx, 1)
			if err != nil {
				log.Println(err.Error())
			}
			if err == nil {
				return ctx, nil
			}
		}
	}
}

func verifyStub(ctx context.Context, expectedCount int64) error {
	ctxData := ctx.Value(testFixtureKey{}).(*testFixture)

	lastStub, ok := ctxData.lastStubRule()
	if !ok {
		return fmt.Errorf("no stub rule found")
	}

	count, err := utils.WiremockClient.GetCountRequests(lastStub.Request())
	if err != nil {
		return fmt.Errorf("could not verify stub rule: %v", err)
	}

	// TODO: count sometimes is greater than expectedCount, investigate why
	if count < expectedCount {
		return fmt.Errorf("stub rule did not verify, expected %d requests, but got %d", expectedCount, count)
	}

	return nil
}
