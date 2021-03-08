package internal

import (
	"github.com/riposo/riposo/pkg/riposo"
	"github.com/riposo/riposo/pkg/schema"
)

// extra contains an account payload.
type extra struct {
	Password *string `json:"password"`
}

// parseExtra parses extra payload.
func parseExtra(obj *schema.Object) (*extra, error) {
	var p *extra
	if err := obj.DecodeExtra(&p); err != nil {
		return nil, schema.BadRequest(err)
	}
	return p, nil
}

// GetPassword is a password reader.
func (p *extra) GetPassword() string {
	if p.Password != nil {
		return *p.Password
	}
	return ""
}

func (p *extra) hashPassword(hlp *riposo.Helpers) error {
	pass := p.GetPassword()
	if pass == "" {
		return schema.InvalidBody("data.password", "Required")
	}

	hashed, err := hlp.SlowHash(pass)
	if err != nil {
		return schema.InternalError(err)
	}

	p.Password = &hashed
	return nil
}
