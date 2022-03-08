package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func NewClient(host, port string) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	status := client.Ping(context.Background())

	if status.Err() != nil {
		return nil, status.Err()
	}

	return client, nil
}
