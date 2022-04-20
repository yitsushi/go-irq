package main

import (
	"github.com/sirupsen/logrus"
	"github.com/yitsushi/go-irq"
)

func main() {
	logger := logrus.New()
	core := irq.NewCore(10, irq.WithLogger(logger))

	core.Boot(newApplication(logger))
}
