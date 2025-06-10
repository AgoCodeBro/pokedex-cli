package main

import ("fmt"
        "strings"
        "bufio"
        "os"
        "errors"
				"pokedex/internal/pokeapi"
        )

// Config to hold pev and next page
type config struct {
	next  *string
	prev  *string
}

// Commands
type cliCommand struct {
    name        string
    description string
    callback    func(*config) error
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
}


func main() {
    scanner := bufio.NewScanner(os.Stdin)
    initCommands()

		initialUrl := "https://pokeapi.co/api/v2/location-area/"
	  config := config{next : &initialUrl,
										prev : nil,
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
            val.callback(&config)
        }
    }
}

// Command Callbacks
func commandExit(c *config) error {
    fmt.Printf("Closing the Pokedex... Goodbye!\n")
    os.Exit(0)
    return errors.New("IDK how you got here bruh, failed to exit")
}

func commandHelp(c *config) error {
    fmt.Printf("Welcome to the Pokedex!\n   Usage:\n\n")
    for _, value := range commands {
        fmt.Printf("%v: %v\n", value.name, value.description)
    }
    return nil
}

func commandMap(c *config) error {
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

func commandMapb(c *config) error {
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

// Other Functions
func cleanInput(text string) []string {
    var result []string
    text = strings.ToLower(text)
    result = strings.Fields(text)

    return result
}
