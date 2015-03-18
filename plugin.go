package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/cf/configuration/config_helpers"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/skratchdot/open-golang/open"
)

func fatalIf(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func fatalWithMessageIf(err error, msg string) {
	if err != nil {
		fmt.Println("ERROR:", msg)
		os.Exit(1)
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

type appEnv struct {
	System map[string]interface{} `json:"system_env_json"`
}

type appEnvService struct {
	Name        string                 // name of the service
	Label       string                 // label of the service
	Tags        []string               // tags for the service
	Plan        string                 // plan of the service
	Credentials map[string]interface{} // credentials for the service
}

type appEnvServices map[string][]appEnvService

// KibanaMeAppPlugin is the type for the plugin functions
type KibanaMeAppPlugin struct{}

func main() {
	plugin.Start(&KibanaMeAppPlugin{})
}

// Run is the entry function for a cf CLI plugin
func (c *KibanaMeAppPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if len(args) < 3 {
		cliConnection.CliCommand(args[0], "-h")
	}

	if args[0] == "kibana-me-logs" {
		kibanaAppName, appName := args[1], args[2]

		kibanaAppOutput, err := cliConnection.CliCommandWithoutTerminalOutput("app", kibanaAppName)
		fatalWithMessageIf(err, "kibana app does not exist in this org/space")
		kibanaGUID := findAppGUID(cliConnection, kibanaAppName)

		_, err = cliConnection.CliCommandWithoutTerminalOutput("app", appName)
		fatalWithMessageIf(err, "app does not exist in this org/space")
		appGUID := findAppGUID(cliConnection, appName)

		kibanaLogstash, err := findService(cliConnection, kibanaGUID, "logstash14")
		fatalIf(err)
		appLogstash, err := findService(cliConnection, appGUID, "logstash14")
		fatalIf(err)

		if appLogstash.Name != kibanaLogstash.Name {
			fatalIf(fmt.Errorf("app and kibana do not share the same logstash14 service"))
		}

		kibanaURLs, err := getURLFromOutput(kibanaAppOutput)
		fatalIf(err)
		kibanaBaseURL := kibanaURLs[0]
		appURL := fmt.Sprintf("%s/#/dashboard/file/app-logs-%s.json", kibanaBaseURL, appGUID)
		open.Run(appURL)
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
					Usage: "kibana-me-logs <kibana-app-name> <app-name>",
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

func findService(cliConnection plugin.CliConnection, appGUID string, serviceName string) (logstash appEnvService, err error) {
	appQuery := fmt.Sprintf("/v2/apps/%v/env", appGUID)
	cmd := []string{"curl", appQuery}
	output, _ := cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	appEnvs := appEnv{}
	json.Unmarshal([]byte(output[0]), &appEnvs)
	str, err := json.Marshal(appEnvs.System["VCAP_SERVICES"])
	if err != nil {
		return
	}
	services := appEnvServices{}
	json.Unmarshal([]byte(str), &services)
	if len(services[serviceName]) == 0 {
		err = fmt.Errorf("app is not bound to a %s service", serviceName)
		return
	}
	logstash = services[serviceName][0]
	return
}

// extracted from cf-plugin-open
func getURLFromOutput(output []string) ([]string, error) {
	urls := []string{}
	for _, line := range output {
		splitLine := strings.Split(strings.TrimSpace(line), " ")
		if splitLine[0] == "urls:" {
			if len(splitLine) > 1 {
				for p := 1; p < len(splitLine); p++ {
					url := "http://" + strings.Trim(splitLine[p], ",")
					url = strings.TrimSpace(url)
					urls = append(urls, url)
				}

			} else if len(splitLine) == 1 {
				return []string{""}, errors.New("App has no route")
			}
		}
	}
	return urls, nil
}
