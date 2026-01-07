package main

import (
	"context"
	"github.com/Screentime42/gator-go/internal/database"
	"fmt"
	"time"
	"github.com/google/uuid"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	cmd map[string]func(*state, command) error
}


func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.cmd[cmd.Name]

	if ok {
		err := handler(s, cmd)
		if err != nil {
			return fmt.Errorf("Command failed: %w", err)
		} 
		return nil
	}
	return fmt.Errorf("Command not found")
}


func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd[name] = f
}

// handlers hooked up to SQL queries
func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no username input")
	}
	
	if len(cmd.Args) > 1 {
		return fmt.Errorf("more than one username input - command only accepts one username")
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("user: %q doesn't exist", cmd.Args[0])
	}
	
	err = s.cfg.SetUser(cmd.Args[0]) 
	if err != nil {
			return fmt.Errorf("set user unsuccessful: %w", err)
	}
	
	fmt.Println("set user successful!")
	fmt.Printf("current user from state: %s\n", s.cfg.CurrentUserName)
	return nil
	
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no username found")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("too many arguments")
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		return fmt.Errorf("user: %q already exists", cmd.Args[0])
	}

	params := database.CreateUserParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		 	cmd.Args[0],
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("user creation unsucessful: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("set user unsuccessful: %w", err)
	}

	fmt.Printf("user: %q created successfully\n", user.Name)
	return nil
}


func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("db reset unsuccessful: %w", err)
	}

	fmt.Println("db reset successful")
	return nil
}


func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("get users unsuccessful: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
	}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	_ = cmd // use later

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil
}