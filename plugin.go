package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cloudfoundry-community/cf-plugin-kibana-me-logs/cftype"
	"github.com/cloudfoundry/cli/cf/api/resources"
	"github.com/cloudfoundry/cli/cf/configuration/config_helpers"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/elgs/gostrgen"
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
	cliConnection  plugin.CliConnection
	shouldAuth     bool
	kibanaUser     string
	kibanaPassword string
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
	c.shouldAuth = true
	if len(args) > 2 && args[2] == "--no-auth" {
		c.shouldAuth = false
	}

	confRepo := core_config.NewRepositoryFromFilepath(config_helpers.DefaultFilePath(), fatalIf)
	spaceGUID := confRepo.SpaceFields().Guid

	_, err := cliConnection.CliCommandWithoutTerminalOutput("app", appName)
	fatalWithMessageIf(err, "app does not exist in this org/space")
	appGUID := c.findAppGUID(spaceGUID, appName)

	logstashGUID, logstashName, err := c.findServiceInstanceGUIDName(appGUID, "logstash14")
	if err != nil {
		fatalIf(fmt.Errorf("App `%s' is not draining logs to a logstash14 service.", appName))
	}

	boundApps, err := c.findAppsBoundToService(logstashGUID)
	kibana, err := c.filterAppWithStartCommand(boundApps, "kibana-me-logs")
	if err != nil {
		fmt.Printf("App `%s' service `%s' not yet bound to a kibana-me-logs app. Deploying...\n", appName, logstashName)
		err = c.cloneAndDeployKibanaMeLogs(logstashName)
		fatalIf(err)
	}

	boundApps, err = c.findAppsBoundToService(logstashGUID)
	kibana, err = c.filterAppWithStartCommand(boundApps, "kibana-me-logs")
	if err != nil {
		fatalIf(fmt.Errorf("Deployment failed. App `%s' service `%s' still not bound to a kibana-me-logs app.", appName, logstashName))
	}

	fullRoute, err := c.firstAppRoute(kibana)
	fatalIf(err)

	fmt.Printf("\n\nDeployment successful!\n")
	if c.shouldAuth {
		fmt.Printf("Username: %s\n", c.kibanaUser)
		fmt.Printf("Password: %s\n", c.kibanaPassword)
	}

	kibanaBaseURL := c.routeToURI(confRepo.IsSSLDisabled(), fullRoute)
	fmt.Printf("\nYou can see your kibana-me-logs app at:\n%s\n", kibanaBaseURL)

	appURL := fmt.Sprintf("%s/#/dashboard/file/app-logs-%s.json", kibanaBaseURL, appGUID)
	open.Run(appURL)
}

// GetMetadata is a CF plugin method for metadata about the plugin
func (c *KibanaMeAppPlugin) GetMetadata() plugin.PluginMetadata {
	versionParts := strings.Split(string(VERSION), ".")
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

func (c *KibanaMeAppPlugin) findServiceInstanceGUIDName(appGUID string, needle string) (serviceInstanceGUID string, serviceInstanceName string, err error) {
	// http://apidocs.cloudfoundry.org/204/apps/list_all_service_bindings_for_the_app.html
	// then find which binding -> maps to service with "serviceName"
	//   -> service_instance_url -> entity.name

	bindings := &cftype.ListAllServiceBindingsForTheApp{}

	cmd := []string{"curl", fmt.Sprintf("/v2/apps/%s/service_bindings", appGUID)}
	output, _ := c.cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &bindings)

	for _, binding := range bindings.Resources {
		serviceInstance, err := c.findServiceInstance(binding.Entity.ServiceInstanceURL)
		fatalIf(err)
		service, err := c.findServiceFromInstance(serviceInstance)
		fatalIf(err)
		for _, tag := range service.Entity.Tags {
			if tag == needle {
				return serviceInstance.Metadata.GUID, serviceInstance.Entity.Name, nil
			}
		}
	}
	return "", "", fmt.Errorf("No service bindings for services tagged %s", needle)
}

func (c *KibanaMeAppPlugin) findServiceInstance(serviceInstanceURL string) (service *cftype.RetrieveAParticularServiceInstance, err error) {
	instance := &cftype.RetrieveAParticularServiceInstance{}

	cmd := []string{"curl", serviceInstanceURL}
	output, err := c.cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	if err != nil {
		return instance, err
	}
	json.Unmarshal([]byte(strings.Join(output, "")), &instance)
	return instance, nil
}

func (c *KibanaMeAppPlugin) findServiceFromInstance(serviceInstance *cftype.RetrieveAParticularServiceInstance) (service *cftype.RetrieveAParticularService, err error) {
	servicePlan := &cftype.RetrieveAParticularServicePlan{}
	service = &cftype.RetrieveAParticularService{}

	cmd := []string{"curl", serviceInstance.Entity.ServicePlanURL}
	output, err := c.cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	if err != nil {
		return service, err
	}
	json.Unmarshal([]byte(strings.Join(output, "")), &servicePlan)

	cmd = []string{"curl", servicePlan.Entity.ServiceURL}
	output, err = c.cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	if err != nil {
		return service, err
	}
	json.Unmarshal([]byte(strings.Join(output, "")), &service)

	return service, nil
}

func (c *KibanaMeAppPlugin) findAppsBoundToService(serviceInstanceGUID string) (appGUIDs []string, err error) {
	bindings := &cftype.ListAllServiceBindingsForTheServiceInstance{}

	appQuery := fmt.Sprintf("/v2/service_instances/%s/service_bindings", serviceInstanceGUID)
	cmd := []string{"curl", appQuery}
	output, _ := c.cliConnection.CliCommandWithoutTerminalOutput(cmd...)

	json.Unmarshal([]byte(strings.Join(output, "")), &bindings)

	for _, binding := range bindings.Resources {
		appGUIDs = append(appGUIDs, binding.Entity.AppGUID)
	}
	return
}

func (c *KibanaMeAppPlugin) filterAppWithStartCommand(appGUIDs []string, startCommand string) (app *cftype.RetrieveAParticularApp, err error) {
	for _, appGUID := range appGUIDs {
		app = &cftype.RetrieveAParticularApp{}
		cmd := []string{"curl", fmt.Sprintf("/v2/apps/%s", appGUID)}
		output, _ := c.cliConnection.CliCommandWithoutTerminalOutput(cmd...)
		json.Unmarshal([]byte(strings.Join(output, "")), &app)
		if app.Entity.DetectedStartCommand == startCommand {
			return app, nil
		}
	}
	return nil, fmt.Errorf("No application found with start command '%s'", startCommand)
}

func (c *KibanaMeAppPlugin) firstAppRoute(app *cftype.RetrieveAParticularApp) (fullRoute string, err error) {
	routes := &cftype.ListAllRoutesForTheApp{}
	cmd := []string{"curl", app.Entity.RoutesURL}
	output, _ := c.cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &routes)

	if routes.TotalResults == 0 {
		return "", fmt.Errorf("App '%s' has no routes", app.Entity.Name)
	}
	route := routes.Resources[0]

	domain := &cftype.RetrieveAParticularDomain{}
	cmd = []string{"curl", route.Entity.DomainURL}
	output, _ = c.cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &domain)

	if route.Entity.Host != "" {
		return fmt.Sprintf("%s.%s", route.Entity.Host, domain.Entity.Name), nil
	}
	return domain.Entity.Name, nil
}

func (c *KibanaMeAppPlugin) routeToURI(isSSLDisabled bool, route string) string {
	if c.shouldAuth {
		route = fmt.Sprintf("%s:%s@%s", c.kibanaUser, c.kibanaPassword, route)
	}
	if isSSLDisabled {
		return fmt.Sprintf("http://%s", route)
	}
	return fmt.Sprintf("https://%s", route)
}

// Uses kibanaMeLogsRepo() from distribution_config.go to determine which repo to clone
func (c *KibanaMeAppPlugin) cloneAndDeployKibanaMeLogs(logstashServiceInstanceName string) error {
	kibanaAppName := fmt.Sprintf("kibana-%s", logstashServiceInstanceName)

	var dir string
	if dir = os.Getenv("KIBANA_ME_LOGS_APP_DIR"); dir == "" {
		tmpDir := os.Getenv("TMPDIR")
		if tmpDir == "" {
			tmpDir = "/tmp"
		}
		var err error
		dir, err = ioutil.TempDir(tmpDir, "kibana")
		if err != nil {
			return err
		}

		var repo = kibanaMeLogsRepo()
		fmt.Printf("Cloning %s...\n", repo)
		cmd := exec.Command("git", "clone", repo, dir)
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	fi, err := os.Lstat(dir)
	if err != nil {
		return err
	}

	fmt.Printf("Checking to see if %s is a symlink\n", dir)
	if fi.Mode()&os.ModeSymlink != 0 {
		fmt.Printf("Reading symlink target for %s\n", dir)
		dir, err = os.Readlink(dir)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Using %s as kibana-me-logs source directory\n", dir)
	fmt.Printf("Pushing kibana-me-logs as %s...\n", kibanaAppName)
	cmd := exec.Command("cf", "push", kibanaAppName, "--no-start", "-p", dir)
	err = cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("Binding %s to %s...\n", kibanaAppName, logstashServiceInstanceName)
	cmd = exec.Command("cf", "bind-service", kibanaAppName, logstashServiceInstanceName)
	err = cmd.Run()
	if err != nil {
		return err
	}

	if c.shouldAuth {
		c.kibanaUser = kibanaAppName
		c.kibanaPassword, err = gostrgen.RandGen(20, gostrgen.All, "", "@$#~")
		if err != nil {
			return err
		}

		fmt.Printf("Setting credentials for %s to %s / %s (user / pass)\n", kibanaAppName, c.kibanaUser, c.kibanaPassword)
		cmd = exec.Command("cf", "set-env", kibanaAppName, "KIBANA_USERNAME", c.kibanaUser)
		err = cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("cf", "set-env", kibanaAppName, "KIBANA_PASSWORD", c.kibanaPassword)
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	fmt.Printf("Starting %s...\n", kibanaAppName)
	cmd = exec.Command("cf", "start", kibanaAppName)

	return cmd.Run()
}
