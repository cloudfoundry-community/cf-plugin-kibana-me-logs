package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cloudfoundry/cli/cf/api/resources"
	"github.com/cloudfoundry/cli/cf/configuration/config_helpers"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/plugin"
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

type appEnvService struct {
	Name        string                 // name of the service
	Label       string                 // label of the service
	Tags        []string               // tags for the service
	Plan        string                 // plan of the service
	Credentials map[string]interface{} // credentials for the service
}

type appEnvServices map[string][]appEnvService

// KibanaMeAppPlugin is the type for the plugin functions
type KibanaMeAppPlugin struct {
	cliConnection plugin.CliConnection
}

func main() {
	plugin.Start(&KibanaMeAppPlugin{})
	// fmt.Printf("%#v\n", (&KibanaMeAppPlugin{}).GetMetadata())
}

// Run is the entry function for a cf CLI plugin
func (c *KibanaMeAppPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	c.cliConnection = cliConnection
	if args[0] != "kibana-me-logs" || len(args) < 2 {
		cliConnection.CliCommand(args[0], "-h")
	}

	appName := args[1]

	confRepo := core_config.NewRepositoryFromFilepath(config_helpers.DefaultFilePath(), fatalIf)
	spaceGUID := confRepo.SpaceFields().Guid

	_, err := cliConnection.CliCommandWithoutTerminalOutput("app", appName)
	fatalWithMessageIf(err, "app does not exist in this org/space")
	appGUID := c.findAppGUID(spaceGUID, appName)

	appLogstash, err := c.findService(appGUID, "logstash14")
	fatalIf(err)

	kibana, _ := c.findAppNameBoundToServiceWithStartCommand(appLogstash, "kibana-me-logs")
	fmt.Printf("kibana: %#v\n", kibana)

}

// GetMetadata is a CF plugin method for metadata about the plugin
func (c *KibanaMeAppPlugin) GetMetadata() plugin.PluginMetadata {
	version, err := Asset("VERSION")
	if err != nil {
		fmt.Println("VERSION go-bindata asset not found")
		version = []byte("0.0.0")
	}
	versionParts := strings.Split(string(version), ".")
	major, _ := strconv.Atoi(versionParts[0])
	minor, _ := strconv.Atoi(versionParts[1])
	patch, _ := strconv.Atoi(strings.TrimSpace(versionParts[2]))

	return plugin.PluginMetadata{
		Name: "kibana-me-logs",
		Version: plugin.VersionType{
			Major: major,
			Minor: minor,
			Build: patch,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "kibana-me-logs",
				HelpText: "open kibana-me-logs for an application",

				UsageDetails: plugin.Usage{
					Usage: "kibana-me-logs <app-name>",
				},
			},
		},
	}
}

func (c *KibanaMeAppPlugin) findAppGUID(spaceGUID string, appName string) string {
	appQuery := fmt.Sprintf("/v2/spaces/%v/apps?q=name:%v&inline-relations-depth=1", spaceGUID, appName)
	cmd := []string{"curl", appQuery}

	output, _ := c.cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	res := &resources.PaginatedApplicationResources{}
	json.Unmarshal([]byte(strings.Join(output, "")), &res)

	return res.Resources[0].Resource.Metadata.Guid
}

func (c *KibanaMeAppPlugin) findService(appGUID string, serviceName string) (service appEnvService, err error) {
	appQuery := fmt.Sprintf("/v2/apps/%v/env", appGUID)
	cmd := []string{"curl", appQuery}
	output, _ := c.cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	appEnvs := models.Environment{}
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
	service = services[serviceName][0]
	return
}

func (c *KibanaMeAppPlugin) findAppNameBoundToServiceWithStartCommand(service appEnvService, startCommand string) (appName string, err error) {
	fmt.Printf("service: %#v\n", service)
	return "test", nil
}

// extracted from cf-plugin-open
func (c *KibanaMeAppPlugin) getURLFromOutput(output []string) ([]string, error) {
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
