package main

import (
	"log"
	"os"

	"code.cloudfoundry.org/cf-syslog-cli/internal/command"
	"code.cloudfoundry.org/cli/plugin"
)

type CFSyslogCLI struct{}

func (c CFSyslogCLI) Run(conn plugin.CliConnection, args []string) {
	if len(args) == 0 {
		log.Fatalf("Expected atleast 1 argument, but got 0.")
	}

	switch args[0] {
	case "create-drain":
		command.CreateDrain(conn, args[1:], log.New(os.Stdout, "", 0))
	case "delete-drain":
		command.DeleteDrain(conn, args[1:])
	}
}

func (c CFSyslogCLI) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "CF Syslog CLI Plugin",
		Commands: []plugin.Command{
			{
				Name: "create-drain",
				UsageDetails: plugin.Usage{
					Usage: "create-drain [options] <app-name> <drain-name> <syslog-drain-url>",
					Options: map[string]string{
						"type": "The type of logs to be sent to the syslog drain. Available types: `logs`, `metrics`, and `all`. Default is `logs`",
					},
				},
			},
			{
				Name: "delete-drain",
				UsageDetails: plugin.Usage{
					Usage: "delete-drain <drain-name>",
				},
			},
		},
	}
}

func main() {
	plugin.Start(CFSyslogCLI{})
}
