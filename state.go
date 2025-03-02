package main

import (
	"context"
	"time"
	"fmt"
	"database/sql"
	"Gattor/internal/config"
	"Gattor/internal/database"
	"github.com/google/uuid"
	"log"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

func HandlerLogin(s *State, cmd Command) error {
	if cmd.Args == nil || len(cmd.Args) != 1 {
		return fmt.Errorf("Username is required")
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Given username doesn't exist: %w", err)
		}
		return fmt.Errorf("Failed to get user: %w", err)
	}
	err = s.Cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("The user has been set")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if cmd.Args == nil || len(cmd.Args) != 1 {
		return fmt.Errorf("Username is required")
	}
	username := cmd.Args[0]
	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			newUser, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
				ID: uuid.New(), 
				CreatedAt: time.Now(), 
				UpdatedAt: time.Now(), 
				Name: username,
			})
			if err != nil {
				return fmt.Errorf("Failed to create user: %w", err)
			}
			err = s.Cfg.SetUser(username)
			if err != nil {
				return fmt.Errorf("Failed to set user: %w", err)
			}
			fmt.Println("User was created", newUser)
			return nil
		}
		return fmt.Errorf("Failed to get user: %w", err)
	}
	return fmt.Errorf("user already exists")
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to delete users: %w", err)
	}
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get all users: %w", err)
	}
	currentUser := s.Cfg.CurrentUserName
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Println("* " + user.Name + " (current)")
		} else {
			fmt.Println("* " + user.Name)
		}
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func HandlerAddFeed(s *State, cmd Command, currentUser database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Error argument")
	}

	newFeed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(), 
		CreatedAt: time.Now(), 
		UpdatedAt: time.Now(), 
		Name: cmd.Args[0],
		Url: cmd.Args[1],
		UserID: currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to create feed: %w", err)
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(), 
		UpdatedAt: time.Now(),
		UserID: currentUser.ID,
		FeedID: newFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to create feed follow: %w", err)
	}

	fmt.Println(newFeed)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		username, err := s.Db.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Failed to get username: %w", err)
		}
		fmt.Println(username)
	}

	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Error argument")
	}
	feed, err := s.Db.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Failed to get feed_id: %w", err)
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(), 
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to create feed follow: %w", err)
	}
	fmt.Println(feed.Name)
	fmt.Println(s.Cfg.CurrentUserName)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to get feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Println(feed.Name)
	}
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Error argument")
	}
	err := s.Db.DeleteFollow(context.Background(), database.DeleteFollowParams{
		UserID: user.ID,
		Url: cmd.Args[0],
	})
	if err != nil {
		return fmt.Errorf("Failed to unfollow: %w", err)
	}
	return nil
}

func scrapeFeeds(s *State) {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return 
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.Db, feed)
}

func scrapeFeed(Db *database.Queries, feed database.Feed) {
	err := Db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return 
	}
	RSS, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, elem := range RSS.Channel.Item {
		fmt.Printf("Found post: %s\n", elem.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(RSS.Channel.Item))
}