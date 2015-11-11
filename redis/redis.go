// Package redis
// implements redis pub/sub provider for Goblin
package redis

import (
	"encoding/json"
	"log"
	"time"

	"gopkg.in/redis.v3"
)

// Redis pub/sub provider
type Redis struct {
	client *redis.Client
}

// Init initializes redis connection
func (r *Redis) Init(connString string) error {
	r.client = redis.NewClient(&redis.Options{
		Addr:     connString,
		Password: "",
		DB:       0,
	})

	_, err := r.client.Ping().Result()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// Publish publishes new event to a queue
func (r *Redis) Publish(queue string, what interface{}) error {
	msg, _ := json.Marshal(what)
	err := r.client.Publish(queue, string(msg)).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Subscribe subscribes to a pub/sub queue
func (r *Redis) Subscribe(
	queue string,
	stop chan bool, interval time.Duration) (chan interface{}, error) {

	pubSub, err := r.client.Subscribe(queue)
	if err != nil {
		return nil, err
	}

	c := make(chan interface{})

	go func() {
		t := time.NewTicker(interval)
		for {
			select {
			case <-t.C:
				msg, err := pubSub.Receive()
				if err == nil {
					c <- msg
				}

			case <-stop:
				log.Println("Stop received")
				t.Stop()
				close(c)
				return
			}
		}
	}()

	return c, nil
}
