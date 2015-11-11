/// Package dispatch provides a common functionality
// for goblin json packets sent over pub/sub provider
package dispatch

import (
	"github.com/goblin-ci/dispatch/redis"

	"time"
)

type PubSuber interface {
	Init(string) error
	Publish(string, interface{}) error
	Subscribe(string, chan bool, time.Duration) (chan interface{}, error)
}

// New creates new PubSuber
func NewRedis(connString string) (r PubSuber, err error) {
	r = &redis.Redis{}
	err = r.Init(connString)
	return
}
