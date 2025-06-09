package main

import (
	"fmt"
	"log"
	"os"

	"github.com/KaareSkytte/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var s state
	s.configPointer = &cfg

	var c commands
	c.handlerFunctions = map[string]func(*state, command) error{}

	c.register("login", handlerLogin)

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
