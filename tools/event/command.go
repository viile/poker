package event

const (
	CommandLogin  = "login"
	CommandList   = "list"
	CommandCreate = "create"
	CommandJoin   = "join"
	CommandExit   = "exit"
)

var CommandConfigs = map[string]CommandConfig{
	CommandLogin: {2, 2},
	CommandList:  {1, 1},
	CommandJoin:  {2, 2},
	CommandExit:  {1, 1},
}

type CommandConfig struct {
	MinArgNumbers int
	MaxArgNumbers int
}
