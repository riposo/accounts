package accounts

import (
	"github.com/riposo/accounts/internal"
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/plugin"
)

func init() {
	plugin.Register(plugin.New(
		"accounts",
		map[string]interface{}{
			"description": "Manage user accounts.",
			"url":         "https://github.com/riposo/accounts",
		},
		func(rts *api.Routes) error {
			rts.Resource("/accounts", internal.Model())
			return nil
		},
		nil,
	))
}
