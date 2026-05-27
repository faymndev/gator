package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	config "github.com/faymndev/gator/internal"
	"github.com/faymndev/gator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		panic(err)
	}

	dbQueries := database.New(db)
	ctx := context.Background()

	// register commands
	commands := newCommands()

	commands.register("init", func(s *state, cmd command) error {
		cfg.CurrentUserName = "faymn"
		cfg.DbUrl = "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
		cfg.Write()
		return nil
	})

	commands.register("register", func(s *state, cmd command) error {
		if len(cmd.args) == 0 {
			return errors.New("expected username")
		}

		username := cmd.args[0]
		user, err := s.db.CreateUser(ctx, database.CreateUserParams{
			ID:   uuid.New(),
			Name: username,
		})

		if err != nil {
			return fmt.Errorf("user already exists: %w", err)
		}

		s.cfg.CurrentUserName = user.Name
		s.cfg.Write()

		fmt.Printf("created user %+v\n", user)

		return nil
	})

	commands.register("login", func(s *state, cmd command) error {
		if len(cmd.args) == 0 {
			return errors.New("expected username")
		}

		username := cmd.args[0]
		user, err := s.db.GetUser(ctx, username)
		if err != nil {
			return fmt.Errorf("user does not exist: %w", err)
		}

		s.cfg.CurrentUserName = user.Name
		if err := s.cfg.Write(); err != nil {
			return fmt.Errorf("failed to save login username: %w", err)
		}

		return nil
	})

	commands.register("users", func(s *state, cmd command) error {
		users, err := s.db.GetUsers(ctx)
		if err != nil {
			return fmt.Errorf("failed to get users: %w", err)
		}

		for _, user := range users {
			if user.Name == s.cfg.CurrentUserName {
				fmt.Printf("* %s (current)\n", user.Name)
			} else {
				fmt.Printf("* %s\n", user.Name)
			}
		}

		return nil
	})

	commands.register("reset", func(s *state, cmd command) error {
		err := s.db.Reset(ctx)
		if err != nil {
			return fmt.Errorf("failed to reset database: %w", err)
		}

		return nil
	})

	// execute command
	s := state{cfg: cfg, db: dbQueries}
	err = commands.run(&s, newCommand(os.Args))
	if err != nil {
		panic(err)
	}
}

type state struct {
	cfg *config.Config
	db  *database.Queries
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
