package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

type cliCommand struct {
    name        string
    description string
    callback    func() error
}

var commands = map[string]cliCommand{
			"exit": {
				name: "exit",
				description: "Exit the Pokedex",
				callback: commandExit,
				},
			"help": {
				name: "help",
				description: "Displays a help message",
				callback: commandHelp,
				},
		}


func cleanInput(s string) []string {
	var results []string
	lowerString := strings.ToLower(s)
	delimitedString := strings.TrimSpace(lowerString)
	results = strings.Fields(delimitedString)
	return results
}

func startRepl() {
	userResponse := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		userResponse.Scan()
		resp := cleanInput(userResponse.Text())
		if len(resp) == 0 {
			continue
		}
		if cmd, ok := commands[resp]; ok == true {
			if err := cmd.callback(); err != nil {
				fmt.Printf("error calling command callback: %v", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}

func commandHelp() error {
	fmt.Println("Usage:")
	for _, cmd := range commands {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}
