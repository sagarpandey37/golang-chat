package utils

/** ---------------- Redis Operations --------------- **/

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func StoreMessages(msg ClientsMeta, rdb *redis.Client) {
	channelKey, err := json.Marshal(msg.ChannelKey)
	data, err := json.Marshal(msg)

	if err != nil {
		panic(err)
	}

	if err := rdb.RPush(ctx, msg.ChannelKey.String(), data).Err(); err != nil {
		panic(err)
	}

	log.Printf("Data of channel id %v stored in Redis", string(channelKey))

	return
}

func FetchMessages(rdb *redis.Client, channelKey string) []string {

	log.Printf("Fetching Previous messages of channel id %v from Redis", string(channelKey))

	chatMessages, err := rdb.LRange(ctx, channelKey, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	return chatMessages

}
