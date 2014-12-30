package service

import (
	"testing"
)

func TestNewFrontendSession(t *testing.T) {
	s := NewSession("1", "server-1", nil, nil)
	s.Set("key1", 1)

	NewFrontendSession(s)
	// a, _ := fs.session.Get("key1")
	// t.Logf("%d", a)

	// s.Set("key1", 1)
	// fs.session.Set("k", v)
}
