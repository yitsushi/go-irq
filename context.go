package irq

import "context"

const (
	requestMemoryKey = "requestMemory"
)

func (core *Core) buildContext() (context.Context, context.CancelFunc) {
	baseCtx, cancel := context.WithCancel(context.Background())
	ctxWithRequestMem := context.WithValue(baseCtx, requestMemoryKey, core.RequestMemorySegment)

	return ctxWithRequestMem, cancel
}

func RequestMemory(ctx context.Context, irq uint) (*Memory, error) {
	fn, ok := ctx.Value(requestMemoryKey).(func(uint) *Memory)
	if !ok {
		return nil, InvalidContextError{Key: requestMemoryKey}
	}

	return fn(irq), nil
}
