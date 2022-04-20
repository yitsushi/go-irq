package irq

import "github.com/sirupsen/logrus"

type WithFunc func(c *Core)

func WithLogger(logger *logrus.Logger) WithFunc {
	return func(c *Core) {
		c.logger = logger
	}
}
