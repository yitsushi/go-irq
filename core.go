package irq

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type MessageBus <-chan uint
type APIFunc func(name string, payload interface{}) (interface{}, error)
type ProcessFunc func(ctx context.Context, data ProcessData)
type AttachProcessFunc func(requestedID uint, process ProcessFunc) error

type Core struct {
	irqChan    chan uint
	processMap map[uint]*Memory
	mutex      sync.Mutex
	logger     *logrus.Logger
}

func NewCore(queueSize int, options ...WithFunc) *Core {
	core := &Core{
		irqChan:    make(chan uint, queueSize),
		processMap: map[uint]*Memory{},
		mutex:      sync.Mutex{},
	}

	for _, fn := range options {
		fn(core)
	}

	if core.logger == nil {
		core.logger = logrus.New()
	}

	return core
}

func (core *Core) Boot(app App) {
	ctx, cancel := core.buildContext()
	defer cancel()

	var wg sync.WaitGroup

	app.Boot(func(requestedID uint, process ProcessFunc) error {
		return core.AttachProcess(&wg, ctx, requestedID, process, app.API)
	})
	app.Start(ctx, core.irqChan)

	cancel()

	wg.Wait()

	close(core.irqChan)
}

func (core *Core) RequestMemorySegment(irq uint) *Memory {
	return core.processMap[irq]
}

func (core *Core) AttachProcess(
	wg *sync.WaitGroup,
	ctx context.Context,
	requestedID uint,
	process ProcessFunc,
	api APIFunc,
) error {
	core.logger.Infof("Register IRQ: %d", requestedID)

	if requestedID < 10 {
		return ReservedInterruptRequestIDError{ID: requestedID}
	}

	core.mutex.Lock()
	defer core.mutex.Unlock()

	core.logger.Infof("Find IRQ conflict: %d", requestedID)

	if _, found := core.processMap[requestedID]; found {
		return InterruptRequestIDAlreadyTakenError{ID: requestedID}
	}

	core.processMap[requestedID] = &Memory{}

	interrupt := func() {
		core.irqChan <- requestedID
	}

	core.logger.Infof("Start process with assigned IRQ: %d", requestedID)
	wg.Add(1)

	go func() {
		process(ctx, newProcessData(interrupt, api, core.processMap[requestedID], core.logger))
		core.logger.Infof("Stopped process with assigned IRQ: %d", requestedID)
		wg.Done()
	}()

	return nil
}
