package storage

import (
	"context"
	"math/rand"
	"time"
)

// DoWithBackoff 使用指数退避和抖动执行函数 fn
// - maxRetries: 重试次数（0 表示只尝试一次）
// - baseDelay: 第一次重试前的初始延迟时间
// - jitter: 添加到延迟中的随机抖动
// 该函数会尊重上下文的取消信号
func DoWithBackoff(
	ctx context.Context,
	maxRetries int,
	baseDelay,
	jitter time.Duration,
	fn func() error,
) error {
	// 避免一些故障：可能是由于网络延迟、服务暂时过载、资源竞争等引起的。通过重试，我们可以自动恢复而不需要用户干预。
	var err error
	delay := baseDelay
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			j := time.Duration(rand.Int63n(int64(jitter)))

			// 等待延迟时间 + 抖动时间，或者等待上下文取消
			select {
			case <-time.After(delay + j):
			case <-ctx.Done():
				return ctx.Err()
			}
			// 指数级增加下一次的延迟时间
			delay *= 2
		}
		if err = fn(); err == nil {
			return nil
		}
	}
	return err
}
