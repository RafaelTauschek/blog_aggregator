package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/RafaelTauschek/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	exsits, ok := c.commands[cmd.name]

	if !ok {
		return errors.New("no command found")
	}

	err := exsits(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("no argument provided")
	}

	err := s.cfg.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Println("user has been set")

	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	s := &state{
		cfg: &cfg,
	}

	cmds := &commands{
		commands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not enough arguments provided")
	}

	cmdName := os.Args[1]
	arg := os.Args[2:]

	cmd := command{
		name:      cmdName,
		arguments: arg,
	}

	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
