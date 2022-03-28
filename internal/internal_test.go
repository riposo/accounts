package internal_test

import (
	"testing"

	. "github.com/bsm/ginkgo/v2"
	. "github.com/bsm/gomega"
	"github.com/riposo/accounts/internal"
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/mock"
	"github.com/riposo/riposo/pkg/schema"
)

var _ = Describe("Account Model", func() {
	var subject api.Model
	var txn *api.Txn

	BeforeEach(func() {
		txn = mock.Txn()
		subject = internal.Model{}
	})

	AfterEach(func() {
		Expect(txn.Rollback()).To(Succeed())
	})

	Describe("Create", func() {
		It("validates", func() {
			Expect(subject.Create(txn, "/accounts/*", &schema.Resource{
				Data: &schema.Object{Extra: []byte(`{}`)},
			})).To(MatchError(`data.id in body: Required`))

			Expect(subject.Create(txn, "/accounts/*", &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{}`)},
			})).To(MatchError(`data.password in body: Required`))
		})

		It("hashes password", func() {
			obj := &schema.Object{ID: "alice", Extra: []byte(`{"password":"s3cret"}`)}
			Expect(subject.Create(txn, "/accounts/*", &schema.Resource{Data: obj})).To(Succeed())
			Expect(obj.Extra).To(MatchJSON(`{"password": "$argon2id$v=19$m=32,t=1,p=1$Iw$Lwbwp6a14wc"}`))
		})

		It("makes the account a writer", func() {
			Expect(subject.Create(txn, "/accounts/*", &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{"password":"s3cret"}`)},
			})).To(Succeed())
			Expect(txn.Perms.GetPermissions("/accounts/alice")).To(Equal(schema.PermissionSet{
				"write": {"account:alice"},
			}))
		})
	})

	Describe("Update", func() {
		var obj *schema.Object

		BeforeEach(func() {
			Expect(subject.Create(txn, "/accounts/*", &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{"password":"s3cret"}`)},
			})).To(Succeed())

			var err error
			obj, err = txn.Store.Get("/accounts/alice", true)
			Expect(err).NotTo(HaveOccurred())
		})

		It("validates", func() {
			Expect(subject.Update(txn, "/accounts/alice", obj, &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{}`)},
			})).To(MatchError(`data.password in body: Required`))
		})

		It("hashes password", func() {
			Expect(subject.Update(txn, "/accounts/alice", obj, &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{"password":"upd@ted"}`)},
			})).To(Succeed())
			Expect(obj.Extra).To(MatchJSON(`{"password":"$argon2id$v=19$m=32,t=1,p=1$Iw$Ef1fxv/3Imo"}`))
		})

		It("retains account as writer", func() {
			Expect(subject.Update(txn, "/accounts/alice", obj, &schema.Resource{
				Data:        &schema.Object{ID: "alice", Extra: []byte(`{"password":"upd@ted"}`)},
				Permissions: schema.PermissionSet{"read": {"account:bob"}},
			})).To(Succeed())
			Expect(txn.Perms.GetPermissions("/accounts/alice")).To(Equal(schema.PermissionSet{
				"read":  {"account:bob"},
				"write": {"account:alice"},
			}))
		})
	})

	Describe("Patch", func() {
		var obj *schema.Object

		BeforeEach(func() {
			Expect(subject.Create(txn, "/accounts/*", &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{"password":"s3cret"}`)},
			})).To(Succeed())

			var err error
			obj, err = txn.Store.Get("/accounts/alice", true)
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not validate", func() {
			orig := string(obj.Extra)
			Expect(subject.Patch(txn, "/accounts/alice", obj, &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{}`)},
			})).To(Succeed())
			Expect(obj.Extra).To(MatchJSON(orig))
		})

		It("hashes password", func() {
			Expect(subject.Patch(txn, "/accounts/alice", obj, &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{"password":"upd@ted"}`)},
			})).To(Succeed())
			Expect(obj.Extra).To(MatchJSON(`{"password":"$argon2id$v=19$m=32,t=1,p=1$Iw$Ef1fxv/3Imo"}`))
		})

		It("retains account as writer", func() {
			Expect(subject.Patch(txn, "/accounts/alice", obj, &schema.Resource{
				Data:        &schema.Object{ID: "alice", Extra: []byte(`{"password":"upd@ted"}`)},
				Permissions: schema.PermissionSet{"read": {"account:bob"}},
			})).To(Succeed())
			Expect(txn.Perms.GetPermissions("/accounts/alice")).To(Equal(schema.PermissionSet{
				"read":  {"account:bob"},
				"write": {"account:alice"},
			}))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal")
}
