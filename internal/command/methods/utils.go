package methods

import (
	"context"
	"fmt"

	"github.com/faymndev/gator/internal/command"
)

func InitConfig(s *command.State, cmd command.Command) error {
	s.Cfg.CurrentUserName = "faymn"
	s.Cfg.DbUrl = "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
	s.Cfg.Write()
	return nil
}

func ResetDatabase(s *command.State, cmd command.Command) error {
	err := s.Db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}
	return nil
}
