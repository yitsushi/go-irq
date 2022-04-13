package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yitsushi/go-irq"
)

func main() {
	logger := logrus.New()
	core := irq.NewCore(10, irq.WithLogger(logger))

	core.Boot(newApplication(logger))
}

type application struct {
	logger     *logrus.Logger
	namespaces NamespaceCache
}

func newApplication(logger *logrus.Logger) *application {
	if logger == nil {
		logger = logrus.New()
	}

	return &application{
		namespaces: NamespaceCache{},
		logger:     logger,
	}
}

func (app *application) Boot(attachProcess irq.AttachProcessFunc) {
	if err := attachProcess(15, NamespaceProcess); err != nil {
		app.logger.Error(err)
	}
}

func (app *application) Start(ctx context.Context, irqChannel irq.MessageBus, requestMem irq.RequestMemSegFunc) {
	app.logger.Info("Start application")

	deadline := time.Now().Add(time.Second * 3)

	for {
		select {
		case <-ctx.Done():
			return
		case irq := <-irqChannel:
			mem := requestMem(irq)

			if irq == 15 {
				app.namespaces = *mem.Data.(*NamespaceCache)

				app.logger.Infof("New namespace cache: %+v", app.namespaces)
			}
		default:
		}

		if time.Now().After(deadline) {
			return
		}
	}
}

type NamespaceCache struct {
	Namespaces map[string][]string
}

func NamespaceProcess(ctx context.Context, interrupt func(), memory *irq.Memory) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Millisecond * 500)
		}

		cache, ok := memory.Data.(*NamespaceCache)
		if !ok {
			memory.Data = &NamespaceCache{Namespaces: map[string][]string{}}

			continue
		}

		cache.Namespaces["Default"] = []string{time.Now().String()}

		interrupt()
	}
}
