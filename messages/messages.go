package messages

import "fmt"

type Message struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}

func (m Message) String() string {
	return fmt.Sprintf("detail: %s, code: %d", m.Detail, m.Code)
}

var (
	ErrUnknown = Message{Code: 4000, Detail: "something went wrong"}
)
