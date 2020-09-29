package event

type Event struct {
	Argv          []string
	Argc          int
}

func NewEvent() *Event {
	return &Event{
		Argv:    make([]string, 0),
	}
}

func (e *Event) Append(m string) {
	e.Argv = append(e.Argv,m)
	e.Argc++
}
