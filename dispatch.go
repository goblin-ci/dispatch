/// Package dispatch provides a common functionality
// for goblin json packets sent over pub/sub provider
package dispatch

import (
	"encoding/json"
	"github.com/goblin-ci/dispatch/redis"
	"io"

	"time"
)

// MessageType represents mq message type
type MessageType int

const (
	BuildCommandMsg MessageType = iota
	CommandOutputMsg
)

// Message mq JSON message
type Message struct {
	Type    MessageType `json:"type"`
	Repo    string      `json:"repo"`
	Payload interface{} `json:"payload"`
}

type PubSuber interface {
	io.Writer
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

// NewJSONMessage creates new JSON encoded Message
func NewJSONMessage(t MessageType, repo string, payload interface{}) []byte {
	msg := Message{
		Type:    t,
		Repo:    repo,
		Payload: payload,
	}
	json, _ := json.Marshal(msg)
	return json
}
