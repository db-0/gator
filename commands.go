package main

import (
	"errors"
	"fmt"

	"github.com/db-0/gator/internal/config"
	"github.com/db-0/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.name)
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("Current username has been set to: %v\n", cmd.args[0])
	return nil
}

func handlerRegiser(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.name)
	}

	err := 
}

func (c *commands) run(s *state, cmd command) error {
	cmdFunc, ok := c.registeredCmds[cmd.name]
	if !ok {
		return errors.New("command not found")
	}

	return cmdFunc(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCmds[name] = f
}
