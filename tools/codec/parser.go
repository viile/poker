package codec

import (
	"github.com/viile/poker/tools/event"
	"strings"
)

type Codec struct {

}

func NewCodec() *Codec {
	return &Codec{

	}
}

func (c *Codec) Decode(buf *[]byte) (e *event.Event, err error) {
	e = event.NewEvent()

	msg := strings.ToLower(strings.TrimSpace(string(*buf)))
	var tmp = 0
	var index = 0
	for index < len(msg) {
		if msg[index] == ' ' {
			e.Append(msg[tmp:index])
			index++
			tmp = index
			continue
		}

		if msg[index] == '\r' || msg[index] == '\n' {
			e.Append(msg[tmp:index])
			return
		}

		index++
	}
	e.Append(msg[tmp:index])
	return
}
