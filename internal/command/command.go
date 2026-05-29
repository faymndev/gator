package command

import (
	config "github.com/faymndev/gator/internal"
	"github.com/faymndev/gator/internal/database"
)

type State struct {
	Cfg *config.Config
	Db  *database.Queries
}
