package accounts

import (
	"github.com/riposo/accounts/internal"
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/plugin"
)

func init() {
	plugin.Register("accounts", func(rts *api.Routes) (plugin.Plugin, error) {
		rts.Resource("/accounts", internal.Model())

		return pin{
			"description": "Manage user accounts.",
			"url":         "https://github.com/riposo/accounts",
		}, nil
	})
}

type pin map[string]interface{}

func (p pin) Meta() map[string]interface{} { return map[string]interface{}(p) }
func (pin) Close() error                   { return nil }
