package main

import (
	"context"
	"github.com/Screentime42/gator-go/internal/database"
	"fmt"
	"time"
	"github.com/google/uuid"
)


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