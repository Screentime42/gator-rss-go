package main

import (
	"fmt"
	"log"
	"os"
	"github.com/Screentime42/gator-go/internal/config"
)

type state struct {
	cfg 	*config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}


	programState := state{cfg}
	cmds := commands{
		cmd: make(map[string]func(*state, command)error),
	}


	cmds.register("login", handlerLogin)


	args := os.Args
	if len(args) < 2 {
		fmt.Println("No command entered")
		os.Exit(1)
	}

	newCommand := command{Name: args[1], Args: args[2:]}
	err = cmds.run(&programState, newCommand)
	if err != nil {
		log.Fatal(err)
	}
}