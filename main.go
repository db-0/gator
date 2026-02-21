package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/db-0/gator/internal/config"
	"github.com/db-0/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading configuration file: %v", err)
	}
	fmt.Printf("Read configuration file: %+v\n", cfg)

	db, err := sql.Open("postgres", cfg.DBURL)
	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	programCmds := commands{
		registeredCmds: make(map[string]func(*state, command) error),
	}

	programCmds.register("login", handlerLogin)
	programCmds.register("register", handlerRegister)

	if len(os.Args) < 2 {
		log.Fatalf("not enough arguments")
	}

	programArgs := os.Args[1:]

	cmd := command{
		name: programArgs[0],
		args: programArgs[1:],
	}

	err = programCmds.run(programState, cmd)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}

}
