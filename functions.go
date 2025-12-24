package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"io"
	"encoding/json"
	"net/http"
	"github.com/Lynn-Xy/pokedex/internal/pokecache"
	"time"
	"math/rand"
)

type cliCommand struct {
    name string
    description string
    callback func(c *config, parameters []string) error
}

type config struct {
	next string
	previous string
	cache *pokecache.Cache
	pokedex *Pokedex
}

type jsonResponse struct {
	Count int `json:"count"`
	Next string	`json:"next"`
	Previous string	`json:"previous"`
	Results []map[string]interface{} `json:"results"`
}

type locationAreaResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	GameIndex int `json:"game_index"`
	EncounterMethodRates []map[string]interface{} `json:"encounter_method_rates"`
	Location map[string]interface{} `json:"location"`
	Names []map[string]interface{} `json:"names"`
	PokemonEncounters []map[string]interface{} `json:"pokemon_encounters"`
}

type pokemon struct {
	Id int `json:"id"`
	Name string `json:"name"`
	BaseExperience int `json:"base_experience"`
	Height int `json:"height"`
	IsDefault bool `json:"is_default"`
	Order int `json:"order"`
	Weight int `json:"weight"`
	Abilities []map[string]interface{} `json:"abilities"`
	Forms []map[string]interface{} `json:"forms"`
	GameIndices []map[string]interface{} `json:"game_indices"`
	HeldItems []map[string]interface{} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves []map[string]interface{} `json:"moves"`
	PastTypes []map[string]interface{} `json:"past_types"`
	PastAbilities []map[string]interface{} `json:"past_abilities"`
	Sprites map[string]interface{} `json:"sprites"`
	Cries map[string]interface{} `json:"cries"`
	Species map[string]interface{} `json:"species"`
	Stats []map[string]interface{} `json:"stats"`
	Types []map[string]interface{} `json:"types"`
}

type Pokedex struct {
	Pokemon map[string]*pokemon
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
		pokedex: &Pokedex{
			Pokemon: map[string]*pokemon{},
		},
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
			err := cmd.callback(c, resps[1:])
			if err != nil {
				fmt.Printf("error executing command %s: %v\n", commandName, err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit(c *config, parameters []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, parameters []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(c *config, parameters []string) error {
	var resp []byte
	if cachedResp, found := c.cache.Get(c.next); found == true {
		resp = cachedResp
	} else {
		respHttp, err := newHttpRequest(c.next)
		if err != nil {
			return fmt.Errorf("error making http request to location-areas: %v", err)
		}
		resp, err = httpRespToBytes(respHttp)
		if err != nil {
			return fmt.Errorf("error reading http response body to bytes: %v", err)
		}
		c.cache.Add(c.next, resp)
	}
	respStruct, err2 := bytesToJsonStruct(resp)
	if err2 != nil {
		return fmt.Errorf("error reading bytes to json response: %v", err2)
	}
	err3 := logJsonResp(respStruct)
	if err3 != nil {
		return fmt.Errorf("error logging http response: %v", err3)
	}
	c.next = respStruct.Next
	c.previous = respStruct.Previous
	return nil
}

func commandMapb(c *config, parameters []string) error {
	if c.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	var resp []byte
	if cachedResp, found := c.cache.Get(c.previous); found == true {
		resp = cachedResp
	} else {
		respHttp, err := newHttpRequest(c.previous)
		if err != nil {
			return fmt.Errorf("error making http request to %v: %v", c.previous, err)
		}
		resp, err = httpRespToBytes(respHttp)
		if err != nil {
			return fmt.Errorf("error reading http response body to bytes: %v", err)
		}
		c.cache.Add(c.previous, resp)
	}
	respStruct, err2 := bytesToJsonStruct(resp)
	if err2 != nil {
		return fmt.Errorf("error reading bytes to json response: %v", err2)
	}
	err3 := logJsonResp(respStruct)
	if err3 != nil {
		return fmt.Errorf("error logging http response: %v", err3)
	}
	c.next = respStruct.Next
	c.previous = respStruct.Previous
	return nil
}

func commandExplore(c *config, parameters []string) error {
	var resp []byte
	if cacheResp, found := c.cache.Get(parameters[0]); found == true {
		resp = cacheResp
	} else {
		respHttp, err := newHttpRequest("https://pokeapi.co/api/v2/location-area/" + parameters[0])
		if err != nil {
			return fmt.Errorf("error making http request to %v: %v", parameters[0], err)
		}
		resp, err = httpRespToBytes(respHttp)
		if err != nil {
			return fmt.Errorf("error reading http response body to bytes: %v", err)
		}
		c.cache.Add(parameters[0], resp)
	}
	respStruct := &locationAreaResponse{}
	err2 := json.Unmarshal(resp, respStruct)
	if err2 != nil {
		return fmt.Errorf("error reading bytes to json response: %v", err2)
	}
	for _, encounter := range respStruct.PokemonEncounters {
		fmt.Println(encounter["pokemon"].(map[string]interface{})["name"])
	}
	return nil
}

func commandCatch(c *config, parameters []string) error {
	resp := &pokemon{}
	var respBytes []byte
	if cacheResp, found := c.cache.Get(parameters[0]); found == true {
		respBytes = cacheResp
	} else {
		respHttp, err := newHttpRequest("https://pokeapi.co/api/v2/pokemon/" + parameters[0])
		if err != nil {
			return fmt.Errorf("error making http request to %v: %v", parameters[0], err)
		}
		respBytes, err = httpRespToBytes(respHttp)
		if err != nil {
			return fmt.Errorf("error reading http response body to bytes: %v", err)
		}
		c.cache.Add(parameters[0], respBytes)
	}
	err2 := json.Unmarshal(respBytes, resp)
	if err2 != nil {
		return fmt.Errorf("error reading bytes to json response: %v", err2)
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", parameters[0])
	caught := rand.Intn(resp.BaseExperience)
	if caught < resp.BaseExperience / 2 {
		fmt.Printf("%s escaped!\n", parameters[0])
	} else {
		fmt.Printf("%s was caught!\n", parameters[0])
		c.pokedex.Pokemon[resp.Name] = resp
	}
	return nil
}

func commandInspect(c *config, parameters []string) error {
	pokemonName := parameters[0]
	if pkmn, found := c.pokedex.Pokemon[pokemonName]; found == true {
		fmt.Printf("Name: %s\n", pkmn.Name)
		fmt.Printf("ID: %d\n", pkmn.Id)
		fmt.Printf("Height: %d\n", pkmn.Height)
		fmt.Printf("Weight: %d\n", pkmn.Weight)
		fmt.Printf("Stats:\n")
		for _, stat := range pkmn.Stats {
			fmt.Printf("  %s: %v\n", stat["stat"].(map[string]interface{})["name"], stat["base_stat"])
		}
		fmt.Printf("Types:\n")
		for _, t := range pkmn.Types {
			fmt.Printf("  %s\n", t["type"].(map[string]interface{})["name"])
		}
	} else {
		fmt.Printf("%s is not in your Pokedex. Catch it first!\n", pokemonName)
	}
	return nil
}

func commandPokedex(c *config, parameters []string) error {
	if len(c.pokedex.Pokemon) == 0 {
		fmt.Println("Your Pokedex is empty. Catch some Pokemon first!")
		return nil
	}
	fmt.Println("Your Pokedex contains the following Pokemon:")
	for name := range c.pokedex.Pokemon {
		fmt.Println(name)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex.",
			callback: commandExit,
			},
		"help": {
			name: "help",
			description: "Displays a help message.",
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
		"explore": {
			name: "explore",
			description: "Explore a specific location area by name.",
			callback: commandExplore,
			},
		"catch": {
			name: "catch",
			description: "Catch a pokemon by name.",
			callback: commandCatch,
			},
		"inspect": {
			name: "inspect",
			description: "Inspect a caught pokemon by name.",
			callback: commandInspect,
			},
		"pokedex": {
			name: "pokedex",
			description: "List all caught pokemon in your pokedex.",
			callback: commandPokedex,
			},
		}
}

func bytesToJsonStruct(data []byte) (*jsonResponse, error) {
	resp := &jsonResponse{}
	err := json.Unmarshal(data, resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json response: %v", err)
	}
	return resp, nil
}

func newHttpRequest(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error creating http GET request: %v", err)
	}
	return resp, nil
}

func httpRespToBytes(r *http.Response) ([]byte, error) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading http request response: %v", err)
	}
	return data, nil
}

func logJsonResp(resp *jsonResponse) error {
	for _, value := range resp.Results {
		fmt.Println(value["name"])
	}
	return nil
}