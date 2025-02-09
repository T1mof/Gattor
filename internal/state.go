package internal

import (
	"Gattor/internal/config"
	"fmt"
)

type State struct {
	Cfg *config.Config
}

func HandlerLogin(s *State, cmd Command) error {
	if cmd.Args == nil || len(cmd.Args) != 1 {
		return fmt.Errorf("Username is required")
	}
	err := s.Cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("The user has been set")
	return nil
}
