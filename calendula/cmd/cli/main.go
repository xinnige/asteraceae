package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	cli "github.com/xinnige/asteraceae/calendula/cli"
)

func main() {
	logfile, err := os.OpenFile("cli.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer func() {
		cerr := logfile.Close()
		if cerr != nil {
			log.Printf("error closing file: %v", cerr)
		}
	}()

	log.SetOutput(logfile)
	log.Println("logging...")

	client := cli.NewCLI()
	mapper := map[string]func(){
		"help": client.MethodHelp,
	}

	if len(os.Args) == 1 || os.Args[1] == "-h" {
		commands := make([]string, len(mapper))
		idx := 0
		for key := range mapper {
			commands[idx] = key
			idx++
		}
		sort.Strings(commands)
		printHelp(commands)
		return
	}

	inputCmd := os.Args[1]
	if fn, ok := mapper[inputCmd]; ok {
		fn()
	} else {
		fmt.Printf("subcommand invalid: %q\n", os.Args[1])
		os.Exit(2)
	}
}

func printHelp(commands []string) {

	fmt.Println("Usage: cli <subcommand> [<args>]")
	fmt.Println("Subcommands: ")
	for _, key := range commands {
		fmt.Printf("\t%s\n", key)
	}
}
