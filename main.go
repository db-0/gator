package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/db-0/gator/internal/config"
	"github.com/db-0/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading configuration file: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	programCmds := commands{
		registeredCmds: make(map[string]func(*state, command) error),
	}

	programCmds.register("register", handlerRegister)
	programCmds.register("login", handlerLogin)
	programCmds.register("users", handlerUsers)
	programCmds.register("reset", handlerReset)
	programCmds.register("agg", handlerAgg)

	if len(os.Args) < 2 {
		log.Fatalf("Usage: gator <command> [args...]")
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = programCmds.run(programState, cmd)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}

}
