package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"io"
	"encoding/json"
	"net/http"
	"strconv"
)

type cliCommand struct {
    name string
    description string
    callback func(c *config) error
}

type config struct {
	next string
	previous string
	pageOffset int
}

type jsonResponse struct {
	Count int `json:"count"`
	Next string	`json:"next"`
	Previous string	`json:"previous"`
	Results []map[string]interface{} `json:"results"`
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

	c := &config{
		next: "",
		previous: "",
		pageOffset: 0,
	}

	for {

		fmt.Print("Pokedex > ")
		userResponse.Scan()
		resps := cleanInput(userResponse.Text())

		if len(resps) == 0 {
			continue
		}

		commandName := resps[0]
		if cmd, ok := commands[commandName]; ok == true {
			if commandName == "mapb" {
				if c.pageOffset == 0 {
					fmt.Println("you're on the first page")
					continue
				} else {
					c.pageOffset -= 20
				}
			}
			if err := cmd.callback(c); err != nil {
				fmt.Printf("error calling command callback: %v", err)
			}
			if commandName == "map" {
				c.pageOffset += 20
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\n")
	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(c *config) error {
	endpoint := "location-area"
	offsetUrl := "?offset=" + strconv.Itoa(c.pageOffset)
	fullUrl := endpoint + offsetUrl
	resp, err := newHttpRequest(fullUrl)
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
	return map[string]cliCommand{
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
			description: "Displays the next 20 locations from the pokemon world.",
			callback: commandMap,
			},
		"mapb": {
			name: "mapb",
			description: "Displays the previous 20 locations from the pokemon world.",
			callback: commandMap,
			},
	}

}

func newHttpRequest(endpoint string) (*http.Response, error) {
	fullUrl := "https://pokeapi.co/api/v2/" + endpoint
	resp, err := http.Get(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("error creating http GET request: %v", err)
	}
	return resp, nil
}

func logHttpResponse(r *http.Response) error {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading http request response: %v", err)
	}
	resp := &jsonResponse{}
	err2 := json.Unmarshal(data, resp)
	if err2 != nil {
		return fmt.Errorf("error unmarshaling json response: %v", err2)
	}
	for _, value := range resp.Results {
		fmt.Println(value["name"])
	}
	return nil
}
