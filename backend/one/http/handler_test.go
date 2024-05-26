package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kartpop/cruncan/backend/one/database/onerequest"
	"github.com/kartpop/cruncan/backend/pkg/id"
	"github.com/stretchr/testify/assert"
)

func TestHttpHandler(t *testing.T) {
	type scenario struct {
		name                   string
		jsonBody               string
		mockRepo               onerequest.Repository
		mockIdService          id.Service
		mockProducer           Producer
		expectedStatusCode     int
		expectedResponseSubstr string
	}

	scenarios := []scenario{
		{
			name:                   "success",
			jsonBody:               `{"user_id": "test_user"}`,
			mockRepo:               &mockRepo{},
			mockIdService:          &mockIdService{},
			mockProducer:           &mockProducer{},
			expectedStatusCode:     201,
			expectedResponseSubstr: `{"request_id":"123"}`,
		},
		{
			name:                   "error storing in db",
			jsonBody:               `{"user_id": "test_user"}`,
			mockRepo:               &mockRepo{isError: true},
			mockIdService:          &mockIdService{},
			mockProducer:           &mockProducer{},
			expectedStatusCode:     500,
			expectedResponseSubstr: ErrFailedToSaveRequestToDatabase,
		},
		{
			name:                   "error sending message over producer",
			jsonBody:               `{"user_id": "test_user"}`,
			mockRepo:               &mockRepo{},
			mockIdService:          &mockIdService{},
			mockProducer:           &mockProducer{isError: true},
			expectedStatusCode:     500,
			expectedResponseSubstr: ErrFailedToSendKafkaMessage,
		},
		{
			name:                   "bad json",
			jsonBody:               `{"user_id": "test_user"`,
			mockRepo:               &mockRepo{},
			mockIdService:          &mockIdService{},
			mockProducer:           &mockProducer{},
			expectedStatusCode:     400,
			expectedResponseSubstr: ErrFailedToParseOneRequest,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			// Arrange
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			defer server.Close()

			httpHandler := NewHandler(context.Background(), s.mockRepo, s.mockIdService, s.mockProducer)

			// Act
			req, err := http.NewRequest("POST", "/one", bytes.NewBufferString(s.jsonBody))
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			http.HandlerFunc(httpHandler.Post).ServeHTTP(recorder, req)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, s.expectedStatusCode, recorder.Code)
			recoderBody := recorder.Body.String()
			if !strings.Contains(recoderBody, s.expectedResponseSubstr) {
				t.Errorf("expected response to contain %q, got %q", s.expectedResponseSubstr, recoderBody)
			}
		})
	}

}

// mock repo for testing
type mockRepo struct {
	isError bool
}

func (m *mockRepo) Create(ctx context.Context, req *onerequest.OneRequest) error {
	if m.isError {
		return errors.New("error storing in db")
	}
	return nil
}

func (m *mockRepo) Get(ctx context.Context, reqId string) (*onerequest.OneRequest, error) {
	if m.isError {
		return nil, errors.New("error fetching from db")
	}
	return &onerequest.OneRequest{}, nil
}

// mock idService for testing
type mockIdService struct {
}

func (m *mockIdService) GenerateID() string {
	return "123"
}

// mock producer for testing
type mockProducer struct {
	isError bool
}

func (m *mockProducer) SendMessage(ctx context.Context, message []byte) error {
	if m.isError {
		return errors.New("error sending message over producer")
	}
	return nil
}

func (m *mockProducer) Close() {
}
