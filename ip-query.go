package main

import (
	"fmt"
	"github.com/cloudfoundry/cli/plugin"
	"os"
	"strings"
	"net"
)

type IpQuery struct{}

func (c *IpQuery) GetMetadata() plugin.PluginMetadata{
	primaryUsage := "cf ip-query <IP>"
	secondaryUsage := "   Only apps you have access to will be queried."

	flags := make(map[string]string)

	return plugin.PluginMetadata{
		Name: "ip-query",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 0,
		},
        MinCliVersion: plugin.VersionType{
            Major: 6,
            Minor: 12,
            Build: 0,
        },
		Commands: []plugin.Command{
			{
				Name:     "ip-query",
				HelpText: "Find what apps are running at a specific IP address.",
				UsageDetails: plugin.Usage{
					Usage:   strings.Join([]string{primaryUsage, secondaryUsage}, "\n\n"),
					Options: flags,
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(IpQuery))
}

func (c * IpQuery) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "ip-query" {
		if len(args) < 2 {
			fmt.Println("Incorrect Usage. \n\nIP is a required argument")
			os.Exit(1)
		}

		ip := args[1]
		if net.ParseIP(ip) == nil {
			fmt.Println("Invalid IP address")
			os.Exit(1)
		}	
		fmt.Printf("Searching for apps running on %s that you have permissions to access. Depending on the number of apps, this may take some time.\n\n", ip)

		matchCount := 0

		api := NewCliCcApi(cliConnection)
		newApps, _ := api.StartedApps()
		for _, app := range newApps {
			stats, error := api.AppStats(app.Guid)
			if error != nil {
				fmt.Println(error.Error())
			}
			for _, instance := range stats.AppInstanceStats {
				if ip == instance.Host {
					space, error := api.Space(app.SpaceUrl)
					if error != nil {
						fmt.Println("Error: "  + error.Error())
					}
					matchCount++
					fmt.Printf("Organization: %s, Space: %s, App: %s, Instance: %s\n", space.OrgName, space.Name, app.Name, instance.Instance)
				}
			}
		}
		if matchCount == 0 {
			fmt.Println("No apps found")
		}
	}
}