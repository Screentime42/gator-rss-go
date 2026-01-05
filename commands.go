package main

import "fmt"

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
	