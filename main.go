package main

import (
	"Gattor/internal"
	"Gattor/internal/config"
	"fmt"
	"log"
	"os"
	"database/sql"
	"Gattor/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	db, err := sql.Open("postgres", cfg.DBURL)
	dbQueries := database.New(db)

	state := internal.State{Cfg: &cfg, Db: dbQueries}
	commands := internal.Сommands{
		Commands: map[string]func(*internal.State, internal.Command) error{
			"login": internal.HandlerLogin,
			"register": internal.HandlerRegister,
		},
	}

	if len(os.Args) < 2 {
		log.Fatalf("error reading command")
	}
	command := internal.Command{Name: os.Args[1], Args: os.Args[2:]}
	err = commands.Run(&state, command)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}
}
