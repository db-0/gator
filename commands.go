package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/db-0/gator/internal/config"
	"github.com/db-0/gator/internal/database"
	"github.com/google/uuid"
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

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		os.Exit(1)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("Current username has been set to: %v\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.name)
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
	})

	if err != nil {
		os.Exit(1)
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User %s has been created:\n%v\n", cmd.args[0], user)
	return nil

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
