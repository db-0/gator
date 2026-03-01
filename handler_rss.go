package main

import (
	"context"
	"fmt"
	"time"

	"github.com/db-0/gator/internal/database"
	"github.com/google/uuid"
)

const feedURL = "https://www.wagslane.dev/index.xml"

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	// For now we'll print the entire struct as directed.
	fmt.Println(feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <name> <url>", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Feed added successfully")

	return nil
}

func handlerListFeed(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get feeds: %w", err)
	}

	for _, f := range feeds {
		user, err := s.db.GetUserByID(context.Background(), f.UserID)
		if err != nil {
			return fmt.Errorf("Unable to retrieve username of user_id:%s (%w)", f.UserID, err)
		}
		printFeed(f, user)
	}

	return nil
}

// Helper / Informational functions
func printFeed(feed database.Feed, user database.User) {
	fmt.Printf(" * Name: 	%v\n", feed.Name)
	fmt.Printf(" * URL:		%v\n", feed.Url)
	fmt.Printf(" * User: 	%v\n", user.Name)
}
