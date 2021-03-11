package xcontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoroutineLeak(t *testing.T) {
	old := runtime.NumGoroutine()

	{
		ctx := Background()
		ctx, _ = WithCancel(ctx)
		ctx, _ = WithNotify(ctx, Paused)
		ctx.WaitFor()
		ctx = WithResetSignalers(ctx)
		ctx = WithStdContext(ctx, context.Background())
	}
	runtime.GC()
	runtime.Gosched()
	runtime.GC()
	runtime.Gosched()

	stack := make([]byte, 65536)
	n := runtime.Stack(stack, true)
	stack = stack[:n]
	require.Equal(t, old, runtime.NumGoroutine(), fmt.Sprintf("%s", stack))
}
