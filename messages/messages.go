package messages

type Message struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}

var (
	ErrUnknown = Message{Code: 4000, Detail: "something went wrong"}
)
