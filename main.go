package main

import (
	"context"
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

	c.register("register", handlerRegister)
	c.register("login", handlerLogin)
	c.register("reset", handlerReset)
	c.register("users", handlerListUsers)
	c.register("agg", handlerAgg)
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("feeds", handlerListFeeds)
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("following", middlewareLoggedIn(handlerListFeedFollows))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	c.register("browse", middlewareLoggedIn(handlerBrowse))

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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		sqlNullString := sql.NullString{
			String: s.cfg.CurrentUserName,
			Valid:  true,
		}

		user, err := s.db.GetUser(context.Background(), sqlNullString)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
