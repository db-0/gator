package main

import (
	"errors"
	"fmt"

	"github.com/db-0/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the login command requires a username argument")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Current username has been set to: %v", cmd.args[0])
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {

}
