package main

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yitsushi/go-irq"
)

type application struct {
	logger      *logrus.Logger
	namespaces  NamespaceCache
	clusterList ClusterList
	mutex       sync.RWMutex
}

func newApplication(logger *logrus.Logger) *application {
	if logger == nil {
		logger = logrus.New()
	}

	return &application{
		namespaces: NamespaceCache{},
		logger:     logger,
		mutex:      sync.RWMutex{},
	}
}

func (app *application) Boot(attachProcess irq.AttachProcessFunc) {
	if err := attachProcess(namespaceCacheIRQ, NamespaceProcess); err != nil {
		app.logger.Error(err)
	}

	if err := attachProcess(clientPoolIRQ, ClientPoolProcess); err != nil {
		app.logger.Error(err)
	}
}

func (app *application) Start(ctx context.Context, irqChannel irq.MessageBus) {
	app.logger.Info("Start application")

	deadline := time.Now().Add(time.Second * 3)

	for {
		select {
		case <-ctx.Done():
			return
		case irqValue := <-irqChannel:
			mem, err := irq.RequestMemory(ctx, irqValue)
			if err != nil {
				app.logger.Error(err)
			}

			switch irqValue {
			case namespaceCacheIRQ:
				app.mutex.Lock()
				app.namespaces = *mem.Data.(*NamespaceCache)
				app.mutex.Unlock()

				app.logger.Infof("New namespace cache: %+v", app.namespaces)
			case clientPoolIRQ:
				app.mutex.Lock()
				app.clusterList = *mem.Data.(*ClusterList)
				app.mutex.Unlock()
			default:
				app.logger.Errorf("Unknown IRQ: %d", irqValue)

			}
		default:
		}

		if time.Now().After(deadline) {
			return
		}
	}
}
