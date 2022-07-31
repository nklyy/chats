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

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}
