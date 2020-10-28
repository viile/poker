package event

type Event struct {
	Name string
	Argv []string
	Argc int
}

func NewEvent() *Event {
	return &Event{
		Argv: make([]string, 0),
	}
}

func (e *Event) Append(m string) {
	e.Argv = append(e.Argv, m)
	e.Argc++
}

func (e *Event) Match(m string) bool {
	return e.Name == m
}
