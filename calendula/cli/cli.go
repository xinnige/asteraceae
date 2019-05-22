package cli

import (
	"flag"
	"fmt"

	"github.com/xinnige/asteraceae/calendula/awsapi"
	"github.com/xinnige/asteraceae/calendula/utils"
)

// CLI defines cli controller
type CLI struct {
	AWSAPI        *awsapi.AWSAPI
	AWSClient     *awsapi.AWSIface
	SerialAPI     utils.SerialInterface
	ErrorBehavior flag.ErrorHandling
}

// NewCLI return a CLI controller
func NewCLI() *CLI {
	return &CLI{
		AWSAPI: &awsapi.AWSAPI{},
	}
}

// MethodHelp helps to print help messages
func (cli *CLI) MethodHelp() {

	fmt.Printf("Hello\n")
}
