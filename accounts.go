package accounts

import (
	"context"

	"github.com/riposo/accounts/internal"
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/plugin"
	"github.com/riposo/riposo/pkg/riposo"
)

func init() {
	plugin.Register(plugin.New(
		"accounts",
		map[string]interface{}{
			"description": "Manage user accounts.",
			"url":         "https://github.com/riposo/accounts",
		},
		func(_ context.Context, rts *api.Routes, _ riposo.Helpers) error {
			rts.Resource("/accounts", internal.Model{})
			return nil
		},
		nil,
	))
}
