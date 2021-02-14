package utils

/** ---------------- Redis Operations --------------- **/

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func storeInRedis(msg Message, rdb *redis.Client) {
	json, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	if err := rdb.RPush(ctx, "chat_messages", json).Err(); err != nil {
		panic(err)
	}
}
