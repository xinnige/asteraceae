package cli

import (
	"flag"
  "fmt"
  "log"
	"os"

	"github.com/xinnige/asteraceae/calendula/slackapi"
	"github.com/xinnige/asteraceae/calendula/utils"
)

// SlackCLI defines slackcli controller
type SlackCLI struct {
	*CLI
	client *slackapi.Client
	token  string
}

const (
	cmdListLogs = "list-logs"

  envAccessToken = "ACCESS_TOKEN"
  maxlimit = 9999
)

// NewSlackCLI returns a pointer of SlackCLI instance
func NewSlackCLI() *SlackCLI {
	accessToken := utils.GetEnv(envAccessToken, "")
	return &SlackCLI{
		CLI:    NewCLI(),
		client: slackapi.NewClient(accessToken),
		token:  accessToken,
	}
}

// Commands returns available commands
func (cli *SlackCLI) Commands() map[string]func() {
	mapper := map[string]func(){
		cmdListLogs:   cli.methodListLogs,
	}
	return mapper
}

// methodListLogs returns a list of audit logs
func (cli *SlackCLI) methodListLogs() {
	cmd := flag.NewFlagSet(cmdListLogs, cli.ErrorBehavior)
	limit := cmd.Int("limit", maxlimit,
    "specify the number of results to return")
	latest := cmd.Int("latest", 0,
    "specify the timestamp of most recent event to include")
	oldest := cmd.Int("oldest", 0,
    "specify the timestamp of most oldest event to include")
	action := cmd.String("action", "",
    "specify the name of the action to filter results")
	actor := cmd.String("actor", "",
    "specify the user ID to filter results")
	entity := cmd.String("entity", "",
    "specify the ID of the target entity to filter results")

	err := cmd.Parse(os.Args[2:])
	if err != nil || !cmd.Parsed() {
		fmt.Printf("Cannot parse arguments (%v)\n", err)
		return
	}

  entries, err :=  cli.client.ListAuditLogs(
    *limit, *latest, *oldest, *action, *actor, *entity)

  if err != nil {
    log.Printf("MethodListLogs Error: %v", err)
    fmt.Printf("Error: %v\n", err)
  }
  fmt.Printf("Found log entries %d\n", len(entries))
  fmt.Println("----------------------")
  jsonBytes := utils.Marshal(entries, &utils.JSONAPI{})
  fmt.Printf("%s\n", jsonBytes)
}
