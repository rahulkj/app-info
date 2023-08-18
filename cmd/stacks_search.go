package cmd

import (
	"encoding/json"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type Stack struct {
	Entity StackEntity `json:"entity"`
}

type StackEntity struct {
	Name string `json:"name"`
}

func getStacks(app *AppSearchResource, cli plugin.CliConnection) {
	var stack Stack
	cmd := []string{"curl", app.Entity.StackUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &stack)

	app.Entity.Stack = stack.Entity.Name
}
