package internal

import (
	"context"
	"time"
	"fmt"
	"database/sql"
	"Gattor/internal/config"
	"Gattor/internal/database"
	"github.com/google/uuid"
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
	rss, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Failed to fetch: %w", err)
	}
	fmt.Println(rss)
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Error argument")
	}

	currentUser, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Failed to get user id: %w", err)
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