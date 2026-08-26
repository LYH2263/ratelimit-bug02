package ratelimit

import (
	"context"
	"errors"
	"time"
)

func (l *Limiter) Wait(ctx context.Context, key string) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		l.mu.RLock()
		closed := l.closed
		l.mu.RUnlock()
		if closed {
			return ErrClosed
		}
		ok, err := l.AllowCtx(ctx, key)
		if ok {
			return nil
		}
		// 配额耗尽是正常限流，应继续等待令牌补充，而非上送成返回值；
		// 否则网关把限流当作内部错误，无法区分二者。仅内部错误才终止。
		if err != nil && !errors.Is(err, ErrExhausted) {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Millisecond):
		}
	}
}
