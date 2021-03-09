package main

import (
	"github.com/riposo/accounts/internal"
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/plugin"
)

var _ plugin.Factory = Plugin

// Plugin export definition.
func Plugin(rts *api.Routes) (plugin.Plugin, error) {
	rts.Resource("/accounts", internal.Model())

	return plugin.New(
		"accounts",
		map[string]interface{}{
			"description": "Manage user accounts.",
			"url":         "https://github.com/riposo/accounts",
		},
		nil,
	), nil
}

func main() {}
