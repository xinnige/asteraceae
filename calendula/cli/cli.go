package cli

import (
	"flag"

	"github.com/xinnige/asteraceae/calendula/awsapi"
	"github.com/xinnige/asteraceae/calendula/utils"
)

const (
	envToken    = "AUTH_TOKEN"
	envEndpoint = "AUTH_ENDPOINT"
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
