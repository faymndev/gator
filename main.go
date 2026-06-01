package main

import (
	"database/sql"
	"os"

	config "github.com/faymndev/gator/internal"
	"github.com/faymndev/gator/internal/command"
	"github.com/faymndev/gator/internal/command/methods"
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
	commands.Register("init", methods.InitConfig)
	commands.Register("register", methods.RegisterUser)
	commands.Register("login", methods.LoginUser)
	commands.Register("users", methods.ListUsers)
	commands.Register("reset", methods.ResetDatabase)
	commands.Register("agg", methods.Aggregate)
	commands.Register("addfeed", methods.AddFeed)
	commands.Register("feeds", methods.ListFeeds)
	commands.Register("follow", methods.FollowFeed)
	commands.Register("following", methods.ListFollowing)

	// execute command
	state := command.State{Cfg: cfg, Db: database.New(db)}
	err = commands.Run(&state, command.NewCommand(os.Args))
	if err != nil {
		panic(err)
	}
}
