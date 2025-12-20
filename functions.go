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

func cleanInput(s string) []string {
	var results []string
	lowerString := strings.ToLower(s)
	delimitedString := strings.TrimSpace(lowerString)
	results = strings.Fields(delimitedString)
	return results
}

func startRepl() {
	userResponse := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	for {
		fmt.Print("Pokedex > ")
		userResponse.Scan()
		resps := cleanInput(userResponse.Text())
		if len(resps) == 0 {
			continue
		}
		commandName := resps[0]
		if cmd, ok := commands[commandName]; ok == true {
			if err := cmd.callback(); err != nil {
				fmt.Printf("error calling command callback: %v", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\n")
	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
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
	return commands

}
