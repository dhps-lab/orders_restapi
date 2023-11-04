package db

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/dhps-lab/orders_restapi/models"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func RedisClient() {
	redisURL := os.Getenv("REDIS_URL")
	redisPass := os.Getenv("REDIS_PASS")
	Client = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPass,
		DB:       0,
	})
	_, err := Client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Unable to connect to Redis: ", err)
	}
	log.Println("Connected to Redis Server")
}

func Publish_order(order models.WorkOrder) error {
	redisQueue := os.Getenv("REDIS_QUEUE")
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(order); err != nil {
		return err
	}
	err := Client.Publish(context.Background(), redisQueue, b.Bytes()).Err()
	log.Println(err)
	if err != nil {
		panic(err)
	}
	return err

}
