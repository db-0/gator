package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	cmdFunc, ok := c.registeredCmds[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}

	return cmdFunc(s, cmd)
}
