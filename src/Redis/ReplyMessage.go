package Redis

type ReplyType int

const (
	_ ReplyType = iota
	SingleLine
	Error
	Integer
	Bulk
	Unknown
)

type replyMessage struct {
	replyType ReplyType
	message   string
}

type ReplyMessage interface {
	GetType() ReplyType
	GetMessage() string
}

func (item *replyMessage) GetType() ReplyType {
	return item.replyType
}

func (item *replyMessage) GetMessage() string {
	return item.message
}
