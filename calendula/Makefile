PKG_NAME:=.
PKG_DIR:=$(PKG_NAME)

AUTH0CLI_NAME:=auth0cli
SLACKCLI_NAME:=slackcli

AUTH0CLI_SRCS:= $(wildcard $(PKG_DIR)/awsapi/*.go) $(wildcard $(PKG_DIR)/mock/*.go) $(wildcard $(PKG_DIR)/astermisc/*.go)
SLACKCLI_SRCS:= $(wildcard $(PKG_DIR)/slackapi/*.go) $(wildcard $(PKG_DIR)/mock/*.go) $(wildcard $(PKG_DIR)/astermisc/*.go)

PWD = $(shell pwd)
NAME:= $(notdir $(PWD))
# ENV?= dev

AUTH0_CLI_BIN:=./$(AUTH0CLI_NAME)
SLACK_CLI_BIN:=./$(SLACKCLI_NAME)

.PHONY: setup
setup:
	GO111MODULE=on go get golang.org/x/tools/cmd/goimports@v0.0.0-20190305010923-fd53dfa0871c
	GO111MODULE=on go get golang.org/x/tools/cmd/cover@v0.0.0-20190305010923-fd53dfa0871c
	GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint
	GO111MODULE=on go get github.com/axw/gocov/gocov
	GO111MODULE=on go get github.com/AlekSi/gocov-xml
	GO111MODULE=on go get github.com/jstemmer/go-junit-report
	GO111MODULE=on go get github.com/golang/mock/gomock
	GO111MODULE=on go get github.com/stretchr/testify/assert
	GO111MODULE=on go get

.PHONY: fmt
fmt:
	set -x&&find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

.PHONY: lint
lint:
	GO111MODULE=on go get $(PKG_DIR)/...
	GO111MODULE=on golangci-lint run ./...


.PHONY: test
test:
	GO111MODULE=on go test $(PKG_DIR)/...

.PHONY: cover
cover:
	GO111MODULE=on go test -v -coverprofile=coverage.txt -covermode count $(PKG_DIR)/... | go-junit-report >report.xml
	GO111MODULE=on gocov test ./... | gocov-xml > coverage.xml
	find . -name coverage.xml|xargs -n1 sed -i -e 's#github.com/[^/]*/[^/]*/##g'
	find . -name coverage.xml|xargs -n1 sed -i -e 's#/go/src#./#g'

.PHONY: build
build: $(AUTH0_CLI_BIN) $(SLACK_CLI_BIN)

$(AUTH0_CLI_BIN): $(AUTH0CLI_SRCS)
	go build -o $(AUTH0_CLI_BIN) $(PKG_DIR)/cmd/$(AUTH0CLI_NAME)/main.go

$(SLACK_CLI_BIN): $(SLACKCLI_SRCS)
	go build -o $(SLACK_CLI_BIN) $(PKG_DIR)/cmd/$(SLACKCLI_NAME)/main.go


.PHONY: clean
clean:
	find . -name debug.test | xargs rm -f
	find . -name "coverage*" | xargs rm -f
	find . -name "*.xml" | xargs rm -f
	rm -f $(AUTH0_CLI_BIN) $(SLACK_CLI_BIN) cli.log
