package accesstoken

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/avast/retry-go"
)

type clientCache struct {
	client Client
	token  atomic.Pointer[Token]
	now    func() time.Time
	mutex  *sync.Mutex
}

func NewClientCache(client Client, now func () time.Time) *clientCache {
	cl := &clientCache{
		client: client,
		now:    now,
		mutex:  &sync.Mutex{},
	}
	cl.token.Store(&Token{})
	return cl
}

func (c *clientCache) GetToken(ctx context.Context) (*Token, error) {
	loadedToken := c.token.Load()
	if c.isTokenExpired(ctx, loadedToken) {
		return c.getAndCacheToken(ctx)
	}

	return loadedToken, nil
}

func (c *clientCache) isTokenExpired(ctx context.Context, token *Token) bool {
	format := "2006-01-02 15:04:05 -0700 MST"
	expiredDate := token.ExpiresAt
	stringDate := expiredDate.Format(format)

	formattedDate, err := time.Parse(format, stringDate)
	if err != nil {
		// If there is an error parsing the date, we consider the token as expired
		return true
	}

	return formattedDate.Before(c.now())
}

func (c *clientCache) getAndCacheToken(ctx context.Context) (*Token, error) {
	// in case of concurrency, only the first call will pass through
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// when second call reaches here the token is already cached
	loadedToken := c.token.Load()
	if !c.isTokenExpired(ctx, loadedToken) {
		return loadedToken, nil
	}

	var token *Token
	if err := retry.Do(
		func() (err error) {
			token, err = c.client.GetToken(ctx)
			return
		},
		retry.Attempts(5),
		retry.LastErrorOnly(true),
	); err != nil {
		return nil, err
	}
	c.token.Store(token)
	return token, nil
}
