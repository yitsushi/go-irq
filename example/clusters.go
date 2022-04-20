package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/yitsushi/go-irq"
)

type Cluster struct {
	Name   string
	Client int64 // fake
}

type ClusterList struct {
	Pool []Cluster
}

func newClientPool() *ClusterList {
	return &ClusterList{Pool: []Cluster{}}
}

func ClientPoolProcess(ctx context.Context, data irq.ProcessData) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Millisecond * 500)
		}

		memory := data.Memory()

		pool, ok := memory.Data.(*ClusterList)
		if !ok {
			memory.Data = newClientPool()

			continue
		}

		value, err := rand.Prime(rand.Reader, 16)
		if err != nil {
			data.Logger().Error(err)
		}

		newPool := []Cluster{}

		for n := 0; n < 4; n++ {
			cluster := Cluster{
				Name:   fmt.Sprintf("cluster-%d", n),
				Client: value.Int64(),
			}
			newPool = append(newPool, cluster)
		}

		pool.Pool = newPool

		data.Interrupt()
	}
}
