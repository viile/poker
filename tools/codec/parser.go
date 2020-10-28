package codec

import (
	"github.com/viile/poker/tools/errors"
	"github.com/viile/poker/tools/event"
	"strings"
)

type Codec struct {
}

func NewCodec() *Codec {
	return &Codec{}
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

	if e.Argc <= 0 || len(e.Argv) <= 0 || e.Argc != len(e.Argv) {
		err = errors.ErrParser.WithMsg("命令长度错误")
		return
	}

	if c, ok := event.CommandConfigs[e.Argv[0]]; ok {
		if c.MinArgNumbers > 0 && e.Argc < c.MinArgNumbers {
			err = errors.ErrParser.WithMsg("命令长度过短")
			return
		}
		if c.MaxArgNumbers > 0 && e.Argc > c.MaxArgNumbers {
			err = errors.ErrParser.WithMsg("命令长度过长")
			return
		}
	}

	e.Name = e.Argv[0]

	return
}
