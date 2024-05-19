package accesstoken

import (
	"context"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_clientCache_GetToken(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	now := time.Now()
	tests := []struct {
		name               string
		clientExpectations func(*MockClient)
		args               args
		want               *Token
		expectedErr        string
		cachedToken        *Token
	}{
		{
			name: "generate a token for first time",
			clientExpectations: func(mc *MockClient) {
				mc.EXPECT().GetToken(gomock.Any()).DoAndReturn(
					func(ctx context.Context) (*Token, error) {
						return &Token{
							AccessToken: "test_token",
							TokenType:   "token_type",
							ExpiresIn:   "899",
							ExpiresAt:   now.Add(time.Minute * 15),
						}, nil
					},
				)
			},
			args: args{
				ctx: context.Background(),
			},
			want: &Token{
				AccessToken: "test_token",
				TokenType:   "token_type",
				ExpiresIn:   "899",
				ExpiresAt:   now.Add(time.Minute * 15),
			},
			cachedToken: nil,
		},
		{
			name: "cached token is expired, get a new token",
			clientExpectations: func(mc *MockClient) {
				mc.EXPECT().GetToken(gomock.Any()).DoAndReturn(
					func(ctx context.Context) (*Token, error) {
						return &Token{
							AccessToken: "test_token",
							TokenType:   "token_type",
							ExpiresIn:   "899",
							ExpiresAt:   now.Add(time.Minute * 15),
						}, nil
					},
				)
			},
			args: args{
				ctx: context.Background(),
			},
			want: &Token{
				AccessToken: "test_token",
				TokenType:   "token_type",
				ExpiresIn:   "899",
				ExpiresAt:   now.Add(time.Minute * 15),
			},
			cachedToken: &Token{
				AccessToken: "test_token",
				TokenType:   "token_type",
				ExpiresIn:   "899",
				ExpiresAt:   now.Add(time.Hour * -2),
			},
		},
		{
			name: "cached token is not expired",
			args: args{
				ctx: context.Background(),
			},
			want: &Token{
				AccessToken: "test_token",
				TokenType:   "token_type",
				ExpiresAt:   now.Add(time.Hour * 2),
			},
			cachedToken: &Token{
				AccessToken: "test_token",
				TokenType:   "token_type",
				ExpiresAt:   now.Add(time.Hour * 2),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewMockClient(gomock.NewController(t))
			if tt.clientExpectations != nil {
				tt.clientExpectations(client)
			}

			cc := NewClientCache(client)
			if tt.cachedToken != nil {
				cc.token.Store(tt.cachedToken)
			}

			got, err := cc.GetToken(tt.args.ctx)
			if err != nil || tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_isTokenExpired(t *testing.T) {
	ctx := context.TODO()
	authClient := NewMockClient(gomock.NewController(t))
	cacheClient := NewClientCache(authClient)

	t.Run("token generated 10 minutes ago (has 5 minutes of validity left), token is valid and must return false", func(t *testing.T) {
		cacheClient.now = func() time.Time {
			return time.Now()
		}
		mockedToken := &Token{
			ExpiresAt: cacheClient.now().Add(time.Minute * 5),
		}
		result := cacheClient.isTokenExpired(ctx, mockedToken)
		assert.False(t, result)
	})

	t.Run("token generated in the actual moment, token is valid and must return false", func(t *testing.T) {
		cacheClient.now = func() time.Time {
			return time.Now()
		}
		mockedToken := &Token{
			ExpiresAt: cacheClient.now().Add(time.Minute * 15),
		}
		result := cacheClient.isTokenExpired(ctx, mockedToken)
		assert.False(t, result)
	})

	t.Run("token generated an hour ago, token should be expired and must return true", func(t *testing.T) {
		cacheClient.now = func() time.Time {
			return time.Now()
		}

		mockedToken := &Token{
			ExpiresAt: cacheClient.now().Add(-time.Hour),
		}
		result := cacheClient.isTokenExpired(ctx, mockedToken)
		assert.True(t, result)
	})

	t.Run("token generated a day ago, token should be expired and must return true", func(t *testing.T) {
		cacheClient.now = func() time.Time {
			return time.Now()
		}

		mockedToken := &Token{
			ExpiresAt: cacheClient.now().Add(time.Hour * -24),
		}
		result := cacheClient.isTokenExpired(ctx, mockedToken)
		assert.True(t, result)
	})
}
