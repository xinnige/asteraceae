package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/xinnige/asteraceae/calendula/auth0api"
	"github.com/xinnige/asteraceae/calendula/utils"
)

// Auth0CLI defines auth0cli controller
type Auth0CLI struct {
	*CLI
	client   *auth0api.Auth0Client
	token    string
	endpoint string
}

const (
	cmdGetUser   = "get-user"
	cmdListUsers = "list-users"
)

// NewAuth0CLI return a CLI controller
func NewAuth0CLI() *Auth0CLI {
	token := utils.GetEnv(envToken, "")
	endpoint := utils.GetEnv(envEndpoint, "")
	return &Auth0CLI{
		CLI:      NewCLI(),
		client:   auth0api.NewAuth0Client(token, endpoint),
		token:    token,
		endpoint: endpoint,
	}
}

func (cli *Auth0CLI) audit() error {
	if cli.token == "" {
		return fmt.Errorf("empty auth token")
	}
	if cli.endpoint == "" {
		return fmt.Errorf("empty auth endpoint")
	}
	return nil
}

// Commands returns available commands
func (cli *Auth0CLI) Commands() map[string]func() {
	mapper := map[string]func(){
		cmdGetUser:   cli.methodGetUser,
		cmdListUsers: cli.methodListUser,
	}
	return mapper
}

// methodGetUser helps to get a user info by user_id
func (cli *Auth0CLI) methodGetUser() {
	if err := cli.audit(); err != nil {
		fmt.Printf("CLIError: %v\n", err)
		return
	}

	cmd := flag.NewFlagSet(cmdGetUser, cli.ErrorBehavior)
	name := cmd.String("name", "", "specify the unique user name")

	err := cmd.Parse(os.Args[2:])
	if err != nil || !cmd.Parsed() {
		fmt.Printf("Cannot parse arguments (%v)\n", err)
		return
	}
	if *name == "" {
		fmt.Println("Invalid input parameters, -h to show help message.")
		fmt.Println("user name cannot be empty")
		return
	}

	user, err := cli.client.GetUserByName(*name)
	if err != nil {
		fmt.Printf("Cannot get user of %s\n%v", *name, err)
		return
	}
	fmt.Printf("User: %s\n%+v\n", *name, *user)
}

// methodListUser helps to list users
func (cli *Auth0CLI) methodListUser() {
	if err := cli.audit(); err != nil {
		fmt.Printf("CLIError: %v\n", err)
		return
	}

	cmd := flag.NewFlagSet(cmdListUsers, cli.ErrorBehavior)
	start := cmd.Int("start", 0, "specify the page number to start (from 0)")
	limit := cmd.Int("limit", -1, "specify the number of users to list")

	err := cmd.Parse(os.Args[2:])
	if err != nil || !cmd.Parsed() {
		fmt.Printf("Cannot parse arguments (%v)\n", err)
		return
	}

	users, err := cli.client.ListUsers(*start, *limit)
	if err != nil {
		fmt.Printf("Cannot list users of (%d,%d)\n%v", *start, *limit, err)
		return
	}
	fmt.Printf("Total: %d users", len(users))
}
