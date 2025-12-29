package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/axiomaticretroactivity/pokedex/internal/pokeapi"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.TrimSpace(strings.ToLower(text)))
}

type config struct {
	client           pokeapi.Client
	next             *string
	previous         *string
	pokemonInventory map[string]pokeapi.PokemonResponse
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays names of map areas in pages of 20 at a time; repeated commands page forward",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Returns to the previous map page",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <area name>",
			description: "Displays names of Pokemon found in <area name>",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon name>",
			description: "Attempts to catch a <pokemon name>. Pokemon with higher base experience will be more difficult to catch!",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon name>",
			description: "Displays statistics of <pokemon name> if one has been caught.",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays a list of all pokemon currently in your pokedex. Pokemon names will be added upon being caught!",
			callback:    commandPokedex,
		},
	}
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *config, args []string) error {
	url := ""
	if cfg.next != nil {
		url = *cfg.next
	}

	resp, err := cfg.client.GetLocations(url)
	if err != nil {
		return err
	}

	for _, r := range resp.Results {
		fmt.Println(r.Name)
	}

	cfg.next = resp.Next
	cfg.previous = resp.Previous

	return nil
}

func commandMapb(cfg *config, args []string) error {
	url := ""
	if cfg.previous != nil {
		url = *cfg.previous
	} else {
		fmt.Println("you're on the first page")
		return nil
	}

	resp, err := cfg.client.GetLocations(url)
	if err != nil {
		return err
	}

	for _, r := range resp.Results {
		fmt.Println(r.Name)
	}

	cfg.next = resp.Next
	cfg.previous = resp.Previous

	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Error: area name must be provided")
	}
	name := args[0]

	resp, err := cfg.client.GetLocationArea(name)
	if err != nil {
		return err
	}

	for _, encounter := range resp.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)

	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Error: Pokemon name must be provided")
	}
	name := args[0]

	resp, err := cfg.client.GetPokemon(name)
	if err != nil {
		return err
	}

	fmt.Println("Throwing a Pokeball at " + name + "...")
	chance := rand.Intn(400)
	if chance >= resp.BaseExperience {
		cfg.pokemonInventory[name] = resp
		fmt.Println(name + " was caught!")
	} else {
		fmt.Println(name + " escaped!")
	}
	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Error: Pokemon name must be provided")
	}
	name := args[0]

	data, ok := cfg.pokemonInventory[name]
	if !ok {
		return fmt.Errorf("Error: Pokemon not found. Try catching one!")
	}

	fmt.Println("Name: " + name)
	fmt.Printf("Height: %d\n", data.Height)
	fmt.Printf("Weight: %d\n", data.Weight)
	fmt.Println("Stats:")
	for i := range data.Stats {
		fmt.Printf("  -%s: %d\n", data.Stats[i].Stat.Name, data.Stats[i].BaseStat)
	}
	fmt.Println("Types:")
	for i := range data.Types {
		fmt.Printf("  -%s\n", data.Types[i].Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args []string) error {
	if len(cfg.pokemonInventory) == 0 {
		return fmt.Errorf("Error: No Pokemon found. Try catching one!")
	}

	fmt.Println("Your Pokedex:")
	for name := range cfg.pokemonInventory {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}

func REPL() {
	cfg := &config{
		client:           pokeapi.NewClient(15*time.Second, 5*time.Second),
		pokemonInventory: map[string]pokeapi.PokemonResponse{},
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		command := cleanedInput[0]
		args := cleanedInput[1:]

		validCommand, ok := getCommands()[command]
		if ok {
			err := validCommand.callback(cfg, args)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
