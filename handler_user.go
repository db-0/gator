package main

import (
	"context"
	"fmt"
	"time"

	"github.com/db-0/gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.Name)
	}

	// Singular argument to go run . register <name>
	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})

	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User %s has been created:\n", user.Name)
	printUser(user)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.Name)
	}

	// Singular argument to the command "go run . login <name>"
	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("Current username has been set to: %v\n", name)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to list users: %w", err)
	}

	for _, u := range users {
		if u.Name == s.cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", u.Name)
		} else {
			fmt.Printf("%s\n", u.Name)
		}
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("unable to reset user table: %w", err)
	}

	fmt.Println("user table has been reset")
	return nil
}

// Helper / Informational functions
func printUser(user database.User) {
	fmt.Printf(" * ID:		%v\n", user.ID)
	fmt.Printf(" * Name: 	%v\n", user.Name)
}
