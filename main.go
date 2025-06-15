package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/KaareSkytte/gator/internal/config"
	"github.com/KaareSkytte/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Error opening config: %v", err)
	}

	dbQueries := database.New(db)

	var s state
	s.db = dbQueries
	s.cfg = &cfg

	var c commands
	c.handlerFunctions = map[string]func(*state, command) error{}

	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)

	if len(os.Args) < 2 {
		fmt.Printf("Invalid input: not enough inputs\n")
		os.Exit(1)
	}

	var com command
	com.name = os.Args[1]
	com.arguments = os.Args[2:]

	err = c.run(&s, com)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}
