package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/KaareSkytte/gator/internal/config"
	"github.com/KaareSkytte/gator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	handlerFunctions map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if c.handlerFunctions[cmd.name] == nil {
		return fmt.Errorf("invalid command\n")
	}

	err := c.handlerFunctions[cmd.name](s, cmd)

	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlerFunctions[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("missing username\n")
	}

	name := cmd.arguments[0]
	_, err := s.db.GetUser(context.Background(), sql.NullString{String: name, Valid: true})
	if err == sql.ErrNoRows {
		return fmt.Errorf("user does not exist")
	}
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	fmt.Printf("Successfully logged in as: %s\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("missing username\n")
	}

	name := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), sql.NullString{String: name, Valid: true})

	if err == nil {
		return fmt.Errorf("User already exists")
	}

	if err != sql.ErrNoRows {
		return fmt.Errorf("unexpected error: %v", err)
	}

	id := uuid.New()
	now := time.Now()

	params := database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      sql.NullString{String: name, Valid: true},
	}

	newUser, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Could not create user: %v", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("could not set user in config: %v", err)
	}
	fmt.Printf("User %s created!\n", name)
	fmt.Printf("New user: %+v\n", newUser)
	return nil
}
