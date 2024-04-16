package id

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthIDShouldBe20CharsLong(t *testing.T) {
	nodeIdSvc, err := NewIDServiceFromIP("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	nodeId := nodeIdSvc.GenerateID()

	assert.Len(t, nodeId, 20, "expected 20 chars but got %v", len(nodeId))
}

func TestNodeIDFromIPShouldReturnInt64OrError(t *testing.T) {

	testCases := []struct {
		name     string
		ip       string
		expected int64
		error    error
	}{
		{
			name:     "valid ip 127.0.0.1 gives nodeId 1",
			ip:       "127.0.0.1",
			expected: 1,
		},
		{
			name:     "valid ip 127.0.0.255 gives nodeId 255",
			ip:       "127.0.0.255",
			expected: 255,
		},
		{
			name:     "valid ip 127.1.255.255 gives nodeId 65535",
			ip:       "127.1.255.255",
			expected: 65535,
		},
		{
			name:     "valid ip 127.1.0.1 gives nodeId 1",
			ip:       "127.1.0.1",
			expected: 1,
		},
		{
			name:     "valid ip 127.1.0.1 gives nodeId 257",
			ip:       "127.1.1.1",
			expected: 257,
		},
		{
			name:     "valid ip 127.0.0.1 gives nodeId 1",
			ip:       "127.0.0.1",
			expected: 1,
		},
		{
			name:     "invalid ip 127.0.0 gives err",
			ip:       "127.0.0",
			expected: 0,
			error:    errors.New("invalid ip address: 127.0.0"),
		},
		{
			name:     "ipv6 ::1 gives err",
			ip:       "::1",
			expected: 0,
			error:    errors.New("not an ipv4 address: ::1"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodeID, err := nodeIDFromIP(tc.ip)
			assert.Equal(t, tc.error, err, "expected %v", nodeID)
			assert.Equal(t, tc.expected, nodeID)
		})
	}
}
