package main

import ("fmt"
		"strings"
		"bufio"
		"os"
		"errors"
		)

// Commands
type cliCommand struct {
	name        string
	description string
	callback    func() error
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
								 
}


func main() {
	scanner := bufio.NewScanner(os.Stdin)
	initCommands()


	for true {
		fmt.Printf("Pokedex > ")
		scanner.Scan()
		rawText := scanner.Text()
		inputs := cleanInput(rawText)
		
		val, ok := commands[inputs[0]]
		if !ok {
			fmt.Printf("Not a valid command")
		} else {
			val.callback()
		}
	}
}

// Command Callbacks
func commandExit() error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return errors.New("IDK how you got here bruh, failed to exit")
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\n	Usage:\n\n")
	for _, value := range commands {
		fmt.Printf("%v: %v\n", value.name, value.description)
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
