package main

import (
	"fmt"
)


func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("No username found!")
	}
	
	if len(cmd.Args) > 1 {
		return fmt.Errorf("More than one username found - command only accepts one username")
	}

	err := s.cfg.SetUser(cmd.Args[0]) 
	if err != nil {
			return fmt.Errorf("Set user unsuccessful: %w", err)
		}
	
	fmt.Println("Set user successful!")
	fmt.Printf("Current user from state: %s\n", s.cfg.CurrentUserName)
	return nil
	
}