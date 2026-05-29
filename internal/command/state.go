package command

import "fmt"

type Command struct {
	Name string
	Args []string
}

func NewCommand(args []string) Command {
	cmd := Command{}
	if len(args) >= 2 {
		cmd.Name = (args)[1]
	}
	if len(args) >= 3 {
		cmd.Args = args[2:]
	}

	return cmd
}

type Commands struct {
	Registry map[string]HandleCommand
}

type HandleCommand func(*State, Command) error

func NewCommands() *Commands {
	return &Commands{
		Registry: make(map[string]HandleCommand),
	}
}

func (c *Commands) Register(name string, handle HandleCommand) {
	c.Registry[name] = handle
}

func (c *Commands) Run(s *State, cmd Command) error {
	handle, ok := c.Registry[cmd.Name]
	if !ok {
		return fmt.Errorf("command %q does not exist", cmd.Name)
	}
	return handle(s, cmd)
}
