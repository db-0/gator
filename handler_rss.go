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

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
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

	// Adding user automatically follows the feed
	err = createFeedFollow(s, user, feed)

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

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	err = createFeedFollow(s, user, feed)
	if err != nil {
		return err
	}

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	followedFeeds, err := s.db.GetFeedFollowsByUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, f := range followedFeeds {
		fmt.Printf("%s\n", f.FeedName)
	}
	return nil
}

// Helper / Informational functions
func printFeed(feed database.Feed, user database.User) {
	fmt.Printf(" * Name: 	%v\n", feed.Name)
	fmt.Printf(" * URL:		%v\n", feed.Url)
	fmt.Printf(" * User: 	%v\n", user.Name)
}

func printFollow(user string, feed string) {
	fmt.Printf("User %s is now following %s.\n", user, feed)
}

func createFeedFollow(s *state, user database.User, feed database.Feed) error {
	_, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create follow feed: %w", err)
	}
	printFollow(user.Name, feed.Name)
	return nil
}
