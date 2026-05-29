package main

import (
	"database/sql"
	"os"

	config "github.com/faymndev/gator/internal"
	"github.com/faymndev/gator/internal/command"
	"github.com/faymndev/gator/internal/database"
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

	commands := command.NewCommands()
	commands.Register("init", command.InitConfig)
	commands.Register("register", command.RegisterUser)
	commands.Register("login", command.LoginUser)
	commands.Register("users", command.ListUsers)
	commands.Register("reset", command.ResetDatabase)

	// execute command
	state := command.State{Cfg: cfg, Db: database.New(db)}
	err = commands.Run(&state, command.NewCommand(os.Args))
	if err != nil {
		panic(err)
	}
}
