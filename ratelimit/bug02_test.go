package ratelimit

import (
	"context"
	"errors"
	"testing"
)

func TestBug02_ExhaustedErrorChain(t *testing.T) {
	l, _ := New(Options{})
	defer l.Close()
	_ = l.SetQuota(context.Background(), "k", Quota{Rate: 1, Burst: 1})
	if !l.Allow("k") {
		t.Fatal("prime")
	}
	_, err := l.AllowCtx(context.Background(), "k")
	if err == nil {
		t.Fatal("expected exhausted")
	}
	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("want ErrExhausted, got %v", err)
	}
}
