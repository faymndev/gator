package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/faymndev/gator/internal/database"
	"github.com/faymndev/gator/internal/feed"
	"github.com/google/uuid"
)

func InitConfig(s *State, cmd Command) error {
	s.Cfg.CurrentUserName = "faymn"
	s.Cfg.DbUrl = "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
	s.Cfg.Write()
	return nil
}

func RegisterUser(s *State, cmd Command) error {
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
	s.Cfg.Write()

	fmt.Printf("created user %+v\n", user)

	return nil
}

func LoginUser(s *State, cmd Command) error {
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

func ListUsers(s *State, cmd Command) error {
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

func ResetDatabase(s *State, cmd Command) error {
	err := s.Db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}
	return nil
}

func Aggregate(s *State, cmd Command) error {
	feed, err := feed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	fmt.Printf("%v", feed)
	return nil
}
