package methods

import (
	"context"
	"errors"
	"fmt"

	"github.com/faymndev/gator/internal/command"
	"github.com/faymndev/gator/internal/database"
	"github.com/google/uuid"
)

func RegisterUser(s *command.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("expected username")
	}

	username := cmd.Args[0]
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:   uuid.New(),
		Name: username,
	})

	if err != nil {
		return fmt.Errorf("user already exists: %w", err)
	}

	s.Cfg.CurrentUserName = user.Name
	if err := s.Cfg.Write(); err != nil {
		return fmt.Errorf("failed to save login username: %w", err)
	}

	fmt.Printf("created user %+v\n", user)

	return nil
}

func LoginUser(s *command.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("expected username")
	}

	username := cmd.Args[0]
	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user does not exist: %w", err)
	}

	s.Cfg.CurrentUserName = user.Name
	if err := s.Cfg.Write(); err != nil {
		return fmt.Errorf("failed to save login username: %w", err)
	}

	return nil
}

func ListUsers(s *command.State, cmd command.Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
