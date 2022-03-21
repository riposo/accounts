package internal

import (
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/riposo"
	"github.com/riposo/riposo/pkg/schema"
)

// Model implements the group model.
type Model struct {
	api.DefaultModel
}

// Create overrides.
func (m Model) Create(txn *api.Txn, path riposo.Path, payload *schema.Resource) error {
	// require ID to be explicitly provided
	if payload.Data.ID == "" {
		return schema.InvalidBody("data.id", "Required")
	}

	// process request
	if err := process(txn, payload, true); err != nil {
		return err
	}

	// include the account itself as a writer
	if payload.Permissions == nil {
		payload.Permissions = make(schema.PermissionSet, 1)
	}
	updatePermissions(payload)

	// perform action
	return m.DefaultModel.Create(txn, path, payload)
}

// Update overrides.
func (m Model) Update(txn *api.Txn, path riposo.Path, exst *schema.Object, payload *schema.Resource) error {
	// process request
	if err := process(txn, payload, true); err != nil {
		return err
	}

	// include the account itself as a writer
	updatePermissions(payload)

	// perform action
	return m.DefaultModel.Update(txn, path, exst, payload)
}

// Patch overrides.
func (m Model) Patch(txn *api.Txn, path riposo.Path, exst *schema.Object, payload *schema.Resource) error {
	// process request
	if err := process(txn, payload, false); err != nil {
		return err
	}

	// include the account itself as a writer
	updatePermissions(payload)

	// perform action
	return m.DefaultModel.Patch(txn, path, exst, payload)
}

func process(txn *api.Txn, payload *schema.Resource, mandatory bool) error {
	// parse payload
	extra, err := parseExtra(payload.Data)
	if err != nil {
		return err
	}

	// hash password
	if mandatory || extra.Password != nil {
		if err := extra.hashPassword(txn.Helpers); err != nil {
			return err
		}
		if err := payload.Data.EncodeExtra(extra); err != nil {
			return schema.InternalError(err)
		}
	}

	return nil
}

func updatePermissions(payload *schema.Resource) {
	if payload.Permissions != nil {
		principal := "account:" + payload.Data.ID
		payload.Permissions.Add("write", principal)
	}
}
