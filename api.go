package main

import (
	"github.com/cloudfoundry/cli/plugin"
	"strings"
	"encoding/json"
)

type StartedApp struct {
	Name		string
	Guid 		string
	SpaceUrl	string
	State 		string
}

type AppStats struct {
	AppInstanceStats []AppInstanceStats
}

type AppInstanceStats struct {
	Instance    string
	Host		string
}

type Space struct {
	Guid  		string
	Name 		string
	OrgName		string
}

type CcApi interface{
	StartedApps() ([]StartedApp, error)
	AppStats(guid string) (AppStats, error)
	Space(space_url string) (Space, error)
}

type CliCcApi struct {
	cliConnection plugin.CliConnection
}

func NewCliCcApi(cliConnection plugin.CliConnection) (api CliCcApi){
	api.cliConnection = cliConnection
	return
}

func (api CliCcApi) StartedApps() ([]StartedApp, error) {
	return startedApps(api.cliConnection, "/v2/apps")
}

func startedApps(cliConnection plugin.CliConnection, url string)([]StartedApp, error){
	appsJson, err := cliConnection.CliCommandWithoutTerminalOutput("curl", url)
	if nil != err {
		return nil, err
	}
	data := strings.Join(appsJson, "\n")
	var appsData map[string]interface{}
 	json.Unmarshal([]byte(data), &appsData)

 	apps := []StartedApp{}
 	for _, app := range appsData["resources"].([]interface{}) {
 		entity := app.(map[string]interface{})["entity"].(map[string]interface{})
 		state := entity["state"].(string)
 		if state == "STARTED" {
 			metadata := app.(map[string]interface{})["metadata"].(map[string]interface{})
 			result := StartedApp{
 				Name: entity["name"].(string),
 				Guid: metadata["guid"].(string),
 				SpaceUrl: entity["space_url"].(string),
 				State: state,
 			}
 			apps = append(apps, result)
 		}
 	} 

 	if nil != appsData["next_url"] {
 		next, _ := startedApps(cliConnection, appsData["next_url"].(string))
 		apps = append(apps, next...)
 	} 

 	return apps, err
}

func (api CliCcApi) AppStats(guid string) (AppStats, error) {
	output, _ := api.cliConnection.CliCommandWithoutTerminalOutput("curl", "/v2/apps/" + guid + "/stats")
	data := strings.Join(output, "\n")
	var instances map[string]map[string]interface{}
	json.Unmarshal([]byte(data), &instances)

	appStats := AppStats{}
	for key, value := range instances {
		if value["state"] == "RUNNING" {
			appStats.AppInstanceStats = append(appStats.AppInstanceStats, AppInstanceStats{
				Instance: key,
				Host: value["stats"].(map[string]interface{})["host"].(string),
				})
		}
	}
	return appStats, nil
}

func (api CliCcApi) Space(space_url string) (Space, error) {
	output, error := api.cliConnection.CliCommandWithoutTerminalOutput("curl", space_url)
	if error != nil {
		return Space{}, error
	}
	data := strings.Join(output, "\n")
	var space map[string]interface{}
	json.Unmarshal([]byte(data), &space)
	entity := space["entity"].(map[string]interface{})
	orgName, error := organizationName(api.cliConnection, entity["organization_url"].(string))
	return Space{
		Name: entity["name"].(string),
		OrgName: orgName,
		}, error
}

func organizationName(cliConnection plugin.CliConnection, org_url string) (string, error) {
	output, error := cliConnection.CliCommandWithoutTerminalOutput("curl", org_url)
	if error != nil {
		return "", error
	}
	data := strings.Join(output, "\n")
	var org map[string]interface{}
	json.Unmarshal([]byte(data), &org)
	entity := org["entity"].(map[string]interface{})
	return entity["name"].(string), error
}