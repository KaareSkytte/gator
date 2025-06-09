package main

import (
	"fmt"

	"github.com/KaareSkytte/gator/internal/config"
)

type state struct {
	configPointer *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	handlerFunctions map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if c.handlerFunctions[cmd.name] == nil {
		return fmt.Errorf("invalid command\n")
	}

	err := c.handlerFunctions[cmd.name](s, cmd)

	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlerFunctions[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("missing username\n")
	}

	name := cmd.arguments[0]
	err := s.configPointer.SetUser(name)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	fmt.Printf("Successfully logged in as: %s\n", name)
	return nil
}
