package irq

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type InterruptRequestIDAlreadyTakenError struct {
	ID uint
}

func (e InterruptRequestIDAlreadyTakenError) Error() string {
	return fmt.Sprintf("irq id is already taken: %d", e.ID)
}

type ReservedInterruptREquestIDError struct {
	ID uint
}

func (e ReservedInterruptREquestIDError) Error() string {
	return fmt.Sprintf("reserved irq id: %d", e.ID)
}

type App interface {
	Boot(fn AttachProcessFunc)
	Start(ctx context.Context, irqChannel MessageBus, fn RequestMemSegFunc)
}

type Memory struct {
	Data interface{}
}

type MessageBus <-chan uint
type ProcessFunc func(ctx context.Context, interrupt func(), memory *Memory)
type AttachProcessFunc func(requestedID uint, process ProcessFunc) error
type RequestMemSegFunc func(irq uint) *Memory
type withFunc func(c *Core)

type Core struct {
	irqChan    chan uint
	processMap map[uint]*Memory
	lock       sync.Mutex
	logger     *logrus.Logger
}

func WithLogger(logger *logrus.Logger) withFunc {
	return func(c *Core) {
		c.logger = logger
	}
}

func NewCore(queueSize int, options ...withFunc) *Core {
	core := &Core{
		irqChan:    make(chan uint, queueSize),
		processMap: map[uint]*Memory{},
		lock:       sync.Mutex{},
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	app.Boot(func(requestedID uint, process ProcessFunc) error {
		return core.AttachProcess(&wg, ctx, requestedID, process)
	})
	app.Start(ctx, core.irqChan, core.RequestMemorySegment)

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
) error {
	core.logger.Infof("Register IRQ: %d", requestedID)

	if requestedID < 10 {
		return ReservedInterruptREquestIDError{ID: requestedID}
	}

	core.lock.Lock()
	defer core.lock.Unlock()

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
		process(ctx, interrupt, core.processMap[requestedID])
		core.logger.Infof("Stopped process with assigned IRQ: %d", requestedID)
		wg.Done()
	}()

	return nil
}
