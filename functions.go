package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"io"
	"encoding/json"
	"net/http"
	"pokedex/internal/pokecache"
	"time"
)

type cliCommand struct {
    name string
    description string
    callback func(c *config) error
}

type config struct {
	next string
	previous string
	cache pokecache.Cache
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
		next: "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		previous: "",
		cache: pokecache.NewCache(5 * 60 * time.Second),
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
			err := cmd.callback(c)
			if err != nil {
				fmt.Printf("error executing command %s: %v\n", commandName, err)
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
	if cachedResp, found := c.cache.Get(c.next); found == true {
		respStruct, err2 := readJsonResponse(cachedResp)
		if err2 != nil {
			return fmt.Errorf("error reading cached http json response: %v", err2)
		}
		err3 := logHttpResponse(respStruct)
		if err3 != nil {
			return fmt.Errorf("error logging http response: %v", err3)
		}
		c.next = respStruct.Next
		c.previous = respStruct.Previous
		return nil
	}
	resp, err := newHttpRequest(c.next)
	if err != nil {
		return fmt.Errorf("error making http request to location-areas")
	}
	c.cache.Add(resp, []byte(c.next))
	respStruct, err2 := readJsonResponse(resp)
	if err2 != nil {
		return fmt.Errorf("error reading http json response: %v", err2)
	}
	c.next = respStruct.Next
	c.previous = respStruct.Previous
	return nil
}

func commandMapb(c *config) error {
	if c.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	if cachedResp, found := c.cache.Get(c.previous); found == true {
		respStruct, err2 := readJsonResponse(cachedResp)
		if err2 != nil {
			return fmt.Errorf("error reading cached http json response: %v", err2)
		}
		err3 := logHttpResponse(respStruct)
		if err3 != nil {
			return fmt.Errorf("error logging http response: %v", err3)
		}
		c.next = respStruct.Next
		c.previous = respStruct.Previous
		return nil
	}
	resp, err := newHttpRequest(c.previous)
	if err != nil {
		return fmt.Errorf("error making http request to location-areas")
	}
	c.cache.Add(resp, []byte(c.previous))
	respStruct, err2 := readJsonResponse(resp)
	if err2 != nil {
		return fmt.Errorf("error reading http json response: %v", err2)
	}
	c.next = respStruct.Next
	c.previous = respStruct.Previous
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
			callback: commandMapb,
			},
	}

}

func newHttpRequest(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error creating http GET request: %v", err)
	}
	return resp, nil
}

func readJsonResponse(r *http.Response) (*jsonResponse, error) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading http request response: %v", err)
	}
	resp := &jsonResponse{}
	err2 := json.Unmarshal(data, resp)
	if err2 != nil {
		return nil, fmt.Errorf("error unmarshaling json response: %v", err2)
	}
	return resp, nil
}

func logHttpResponse(resp *jsonResponse) error {
	for _, value := range resp.Results {
		fmt.Println(value["name"])
	}
	return nil
}