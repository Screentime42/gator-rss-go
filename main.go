package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/Screentime42/gator-go/internal/config"
	"github.com/Screentime42/gator-go/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  	*database.Queries
	cfg 	*config.Config
}

func main() {
	// Read config
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{
		db:	dbQueries,
		cfg:	&cfg,
	}


	cmds := commands{
		cmd: make(map[string]func(*state, command)error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerViewFeeds)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("No command entered")
		os.Exit(1)
	}

	newCommand := command{Name: args[1], Args: args[2:]}
	err = cmds.run(programState, newCommand)
	if err != nil {
		log.Fatal(err)
	}
}