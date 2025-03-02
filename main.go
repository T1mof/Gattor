package main

import (
	"Gattor/internal/config"
	"log"
	"os"
	"database/sql"
	"Gattor/internal/database"
	_ "github.com/lib/pq"
	"context"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	//fmt.Printf("Read config: %+v\n", cfg)

	db, err := sql.Open("postgres", cfg.DBURL)
	dbQueries := database.New(db)

	state := State{Cfg: &cfg, Db: dbQueries}
	commands := Сommands{
		Commands: map[string]func(*State, Command) error{
			"login": HandlerLogin,
			"register": HandlerRegister,
			"reset": HandlerReset,
			"users": HandlerUsers,
			"agg": HandlerAgg,
			"addfeed": middlewareLoggedIn(HandlerAddFeed),
			"feeds": HandlerFeeds,
			"follow": middlewareLoggedIn(HandlerFollow),
			"following": middlewareLoggedIn(HandlerFollowing),
			"unfollow": middlewareLoggedIn(HandlerUnfollow),
		},
	}

	if len(os.Args) < 2 {
		log.Fatalf("error reading command")
	}
	command := Command{Name: os.Args[1], Args: os.Args[2:]}
	err = commands.Run(&state, command)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}
}

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}