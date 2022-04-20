package irq

import "github.com/sirupsen/logrus"

type ProcessData interface {
	Interrupt()
	API(name string, payload interface{}) (interface{}, error)
	Memory() *Memory
	Logger() *logrus.Logger
}

type processData struct {
	interrupt func()
	api       APIFunc
	memory    *Memory
	logger    *logrus.Logger
}

func newProcessData(interrupt func(), api APIFunc, memory *Memory, logger *logrus.Logger) ProcessData {
	return &processData{
		interrupt: interrupt,
		api:       api,
		memory:    memory,
		logger:    logger,
	}
}

func (p *processData) Interrupt() {
	p.interrupt()
}

func (p *processData) API(name string, payload interface{}) (interface{}, error) {
	return p.api(name, payload)
}

func (p *processData) Memory() *Memory {
	return p.memory
}

func (p *processData) Logger() *logrus.Logger {
	return p.logger
}
