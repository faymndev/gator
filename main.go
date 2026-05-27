package main

import (
	"errors"
	"fmt"
	"os"

	config "github.com/faymndev/gator/internal"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	s := state{cfg: cfg}

	commands := newCommands()

	commands.register("init", func(s *state, cmd command) error {
		cfg.CurrentUserName = "faymn"
		cfg.DbUrl = "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
		cfg.Write()
		return nil
	})

	commands.register("login", func(s *state, cmd command) error {
		if len(cmd.args) == 0 {
			return errors.New("expected username")
		}

		s.cfg.CurrentUserName = cmd.args[0]
		if err := s.cfg.Write(); err != nil {
			return fmt.Errorf("failed to save login username: %w", err)
		}

		return nil
	})

	err = commands.run(&s, newCommand(os.Args))
	if err != nil {
		panic(err)
	}
}

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	registry map[string]handleCommand
}

type handleCommand func(*state, command) error

func newCommand(args []string) command {
	cmd := command{}
	if len(args) >= 2 {
		cmd.name = (args)[1]
	}
	if len(args) >= 3 {
		cmd.args = args[2:]
	}

	return cmd
}

func newCommands() *commands {
	return &commands{
		registry: make(map[string]handleCommand),
	}
}

func (c *commands) register(name string, handle handleCommand) {
	c.registry[name] = handle
}

func (c *commands) run(s *state, cmd command) error {
	handle, ok := c.registry[cmd.name]
	if !ok {
		return fmt.Errorf("command %q does not exist", cmd.name)
	}
	return handle(s, cmd)
}
