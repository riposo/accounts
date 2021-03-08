package internal

import (
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/conn/storage"
	"github.com/riposo/riposo/pkg/riposo"
	"github.com/riposo/riposo/pkg/schema"
)

type account struct{ api.Model }

// Model inits a new model
func Model() api.Model {
	return &account{Model: api.StdModel()}
}

// Create overrides.
func (m *account) Create(txn *api.Txn, path riposo.Path, payload *schema.Resource) error {
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
	updatePermissions(txn, payload)

	// perform action
	return m.Model.Create(txn, path, payload)
}

// Update overrides.
func (m *account) Update(txn *api.Txn, path riposo.Path, hs storage.UpdateHandle, payload *schema.Resource) error {
	// process request
	if err := process(txn, payload, true); err != nil {
		return err
	}

	// include the account itself as a writer
	updatePermissions(txn, payload)

	// perform action
	return m.Model.Update(txn, path, hs, payload)
}

// Patch overrides.
func (m *account) Patch(txn *api.Txn, path riposo.Path, hs storage.UpdateHandle, payload *schema.Resource) error {
	// process request
	if err := process(txn, payload, false); err != nil {
		return err
	}

	// include the account itself as a writer
	updatePermissions(txn, payload)

	// perform action
	return m.Model.Patch(txn, path, hs, payload)
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

func updatePermissions(txn *api.Txn, payload *schema.Resource) {
	if payload.Permissions != nil {
		principal := "account:" + payload.Data.ID
		payload.Permissions.Add("write", principal)
	}
}
