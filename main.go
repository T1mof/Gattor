package main

import (
	"Gattor/internal"
	"Gattor/internal/config"
	"fmt"
	"log"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	state := internal.State{Cfg: &cfg}
	commands := internal.Сommands{
		Commands: map[string]func(*internal.State, internal.Command) error{
			"login": internal.HandlerLogin,
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
