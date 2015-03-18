package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/cf/configuration/config_helpers"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/plugin"
)

func fatalIf(err error) {
	if err != nil {
		log.Fatalln("error:", err)
	}
}

type appSearchMetaData struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

type appSearchResources struct {
	Metadata appSearchMetaData `json:"metadata"`
}

type appSearchResults struct {
	Resources []appSearchResources `json:"resources"`
}

// KibanaMeAppPlugin is the type for the plugin functions
type KibanaMeAppPlugin struct{}

func main() {
	plugin.Start(&KibanaMeAppPlugin{})
}

// Run is the entry function for a cf CLI plugin
func (c *KibanaMeAppPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if len(args) < 2 {
		cliConnection.CliCommand(args[0], "-h")
	}

	if args[0] == "kibana-me-logs" {
		_, err := cliConnection.CliCommandWithoutTerminalOutput("app", args[1])
		if err != nil {
			fmt.Fprintln(os.Stdout, "error: app does not exist in this org/space")
			os.Exit(1)
		}
		appName := args[1]
		guid := findAppGUID(cliConnection, appName)
		fmt.Println("App GUID:", guid)
	}

}

// GetMetadata is a CF plugin method for metadata about the plugin
func (c *KibanaMeAppPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "kibana-me-logs",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "kibana-me-logs",
				HelpText: "open kibana-me-logs for an application",

				UsageDetails: plugin.Usage{
					Usage: "kibana-me-logs <appname>",
				},
			},
		},
	}
}

func findAppGUID(cliConnection plugin.CliConnection, appName string) string {

	confRepo := core_config.NewRepositoryFromFilepath(config_helpers.DefaultFilePath(), fatalIf)
	spaceGUID := confRepo.SpaceFields().Guid

	appQuery := fmt.Sprintf("/v2/spaces/%v/apps?q=name:%v&inline-relations-depth=1", spaceGUID, appName)
	cmd := []string{"curl", appQuery}

	output, _ := cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	res := &appSearchResults{}
	json.Unmarshal([]byte(strings.Join(output, "")), &res)

	return res.Resources[0].Metadata.GUID
}
