package irq

import "context"

type App interface {
	Boot(fn AttachProcessFunc)
	Start(ctx context.Context, irqChannel MessageBus)
	API(name string, payload interface{}) (interface{}, error)
}
