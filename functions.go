package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"io"
	"encoding/json"
	"net/http"
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

func commandMap() error {
	resp, err := newHttpRequest("location-areas")
	if err != nil {
		return fmt.Errorf("error making http request to location-areas")
	}
	err2 := logHttpResponse(resp)
	if err2 != nil {
		return fmt.Errorf("error reading http json response: %v", err2)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	return commands := map[string]cliCommand{
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
		"map":	{
			name: "map",
			description: "Displays locations from the pokemon world."
			callback: commandMap,
			},
	}

}

func newHttpRequest(endpoint string) (http.Response, error) {
	fullUrl := "https://pokeapi.co/api/v2" + endpoint
	resp, err := http.GET(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("error creating http GET request: %v", err)
	}
	return resp, nil
}

func logHttpResponse(r http.Response) error {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading http request response: %v", err)
	}
	text, err2 := json.Unmarshal(data)
	if err2 != nil {
		return fmt.Errorf("error unmarshaling json response: %v", err2)
	}
	fmt.Println(text)
	return nil
}
