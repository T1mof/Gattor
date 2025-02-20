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