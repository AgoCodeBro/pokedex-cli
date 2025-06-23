package main

import ("fmt"
        "strings"
        "bufio"
        "os"
        "errors"
				"pokedex/internal/pokeapi"
				"math/rand"
        )

// Config to hold pev and next page
type config struct {
	next  *string
	prev  *string
	pokedex map[string]pokeapi.Pokemon
}

// Commands
type cliCommand struct {
    name        string
    description string
    callback    func(*config, []string) error
}

var commands = make(map[string]cliCommand)

func initCommands() {
    commands["help"] = cliCommand{name: "help",
                                  description : "Displays a help message",
                                  callback :  commandHelp,
                                 }
                                 
    commands["exit"] = cliCommand{name: "exit",
                                  description : "Exit the Pokedex",
                                  callback :  commandExit,
                                 }
                                 
		commands["map"] = cliCommand{name : "map",
																 description : "Displays the names of 20 location areas. Calling it again will display the next 20. (see mapb)",
																 callback : commandMap,
																}
	
		commands["mapb"] = cliCommand{name : "mapb",
																 description : "Displays the previous page of locations",
																 callback : commandMapb,
																}

		commands["explore"] = cliCommand{name : "explore",
																		 description : "Takes the name of a location area as an argument and returns the pokemon found there",
																		 callback : commandExplore,
																		}

		commands["catch"] = cliCommand{name : "catch",
																		 description : "Takes the name of a pokemon as an argument and attempts to catch it",
																		 callback : commandCatch,
																		}
}


func main() {
    scanner := bufio.NewScanner(os.Stdin)
    initCommands()

		initialUrl := "https://pokeapi.co/api/v2/location-area/"
	  config := config{next : &initialUrl,
										prev : nil,
										pokedex: make(map[string]pokeapi.Pokemon),
									 }


    for true {
        fmt.Printf("Pokedex > ")
        scanner.Scan()
        rawText := scanner.Text()
        inputs := cleanInput(rawText)
        val, ok := commands[inputs[0]]
        if !ok {
            fmt.Printf("Not a valid command\n")
        } else {
            val.callback(&config, inputs)
        }
    }
}

// Command Callbacks
func commandExit(c *config, inputs []string) error {
    fmt.Printf("Closing the Pokedex... Goodbye!\n")
    os.Exit(0)
    return errors.New("IDK how you got here bruh, failed to exit")
}

func commandHelp(c *config, inputs []string) error {
    fmt.Printf("Welcome to the Pokedex!\n   Usage:\n\n")
    for _, value := range commands {
        fmt.Printf("%v: %v\n", value.name, value.description)
    }
    return nil
}

func commandMap(c *config, inputs []string) error {
	if c.next == nil {
		fmt.Printf("At final page")
		return nil
	}

	url := *c.next
	

	page, err := pokeapi.GetAPIPage(url)
	if err != nil {
		return err
	}

	c.next = page.Next
	c.prev = page.Previous

	for _, result := range page.Results {
		fmt.Printf("%v\n", result.Name)
	}
	
	return nil

}

func commandMapb(c *config, inputs []string) error {
	if c.prev == nil {
		fmt.Printf("At the first page\n")
		return nil
	}
	url := *c.prev
	
	
	page, err := pokeapi.GetAPIPage(url)
	if err != nil {
		return err
	}

	c.next = page.Next
	c.prev = page.Previous

	for _, result := range page.Results {
		fmt.Printf("%v\n", result.Name)
	}

	return nil
}

func commandExplore(c *config, inputs []string) error {
	if len(inputs) < 2 {
		fmt.Printf("Explore requires an argument\n")
		return nil
	}
	response, err := pokeapi.GetLocationArea(inputs[1])
	if err != nil {
		return err
	}

	for _, encounter := range response.PokemonEncounters {
		fmt.Printf("%v\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(c *config, inputs []string) error {
	if len(inputs) < 2 {
		fmt.Printf("Catch requires an argument\n")
		return nil
	}

	response, err := pokeapi.GetPokemon(inputs[1])
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", inputs[1])

	baseLevel := response.BaseExperience

	if roll := rand.Intn(1000); roll >= baseLevel {
		fmt.Printf("%v was caught!\n", inputs[1])
		c.pokedex[inputs[1]] = response
	} else {
		fmt.Printf("%v escaped!\n", inputs[1])
	}

	return nil
}	

// Other Functions
func cleanInput(text string) []string {
    var result []string
    text = strings.ToLower(text)
    result = strings.Fields(text)

    return result
}
