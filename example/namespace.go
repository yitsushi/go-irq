package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/yitsushi/go-irq"
)

type NamespaceCache struct {
	Namespaces map[string][]string
}

func NamespaceProcess(ctx context.Context, data irq.ProcessData) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Millisecond * 200)
		}

		memory := data.Memory()

		cache, ok := memory.Data.(*NamespaceCache)
		if !ok {
			memory.Data = &NamespaceCache{Namespaces: map[string][]string{}}

			continue
		}

		response, err := data.API(requestClientPool, nil)
		if err != nil {
			data.Logger().Error(err)

			continue
		}

		clientPool, ok := response.([]Cluster)
		if !ok || len(clientPool) == 0 {
			continue
		}

		newNamespaceList := map[string][]string{}

		for _, cluster := range clientPool {
			value, err := rand.Prime(rand.Reader, 16)
			if err != nil {
				data.Logger().Error(err)

				continue
			}

			newNamespaceList[cluster.Name] = []string{
				fmt.Sprintf("ns-%d-%d", cluster.Client, value.Int64()),
			}
		}

		cache.Namespaces = newNamespaceList

		data.Interrupt()
	}
}
