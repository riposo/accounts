package internal_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/riposo/accounts/internal"
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/conn/storage"
	"github.com/riposo/riposo/pkg/mock"
	"github.com/riposo/riposo/pkg/schema"
)

var _ = Describe("Account Model", func() {
	var subject api.Model
	var txn *api.Txn

	BeforeEach(func() {
		txn = mock.Txn()
		subject = internal.Model()
	})

	AfterEach(func() {
		Expect(txn.Abort()).To(Succeed())
	})

	Describe("Create", func() {
		It("should validate", func() {
			Expect(subject.Create(txn, "/accounts/*", &schema.Resource{
				Data: &schema.Object{Extra: []byte(`{}`)},
			})).To(MatchError(`data.id in body: Required`))

			Expect(subject.Create(txn, "/accounts/*", &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{}`)},
			})).To(MatchError(`data.password in body: Required`))
		})

		It("should hash password", func() {
			obj := &schema.Object{ID: "alice", Extra: []byte(`{"password":"s3cret"}`)}
			Expect(subject.Create(txn, "/accounts/*", &schema.Resource{Data: obj})).To(Succeed())
			Expect(obj.Extra).To(MatchJSON(`{"password": "$argon2id$v=19$m=32,t=1,p=1$Iw$Lwbwp6a14wc"}`))
		})

		It("should make the account a writer", func() {
			Expect(subject.Create(txn, "/accounts/*", &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{"password":"s3cret"}`)},
			})).To(Succeed())
			Expect(txn.Perms.GetPermissions("/accounts/alice")).To(Equal(schema.PermissionSet{
				"write": {"account:alice"},
			}))
		})
	})

	Describe("Update", func() {
		var hs storage.UpdateHandle

		BeforeEach(func() {
			obj := &schema.Object{ID: "alice", Extra: []byte(`{"password":"s3cret"}`)}
			Expect(subject.Create(txn, "/accounts/*", &schema.Resource{Data: obj})).To(Succeed())

			var err error
			hs, err = txn.Store.GetForUpdate("/accounts/alice")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should validate", func() {
			Expect(subject.Update(txn, "/accounts/alice", hs, &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{}`)},
			})).To(MatchError(`data.password in body: Required`))
		})

		It("should hash password", func() {
			Expect(subject.Update(txn, "/accounts/alice", hs, &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{"password":"upd@ted"}`)},
			})).To(Succeed())
			Expect(hs.Object().Extra).To(MatchJSON(`{"password":"$argon2id$v=19$m=32,t=1,p=1$Iw$Ef1fxv/3Imo"}`))
		})

		It("should retain account as writer", func() {
			Expect(subject.Update(txn, "/accounts/alice", hs, &schema.Resource{
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
		var hs storage.UpdateHandle

		BeforeEach(func() {
			obj := &schema.Object{ID: "alice", Extra: []byte(`{"password":"s3cret"}`)}
			Expect(subject.Create(txn, "/accounts/*", &schema.Resource{Data: obj})).To(Succeed())

			var err error
			hs, err = txn.Store.GetForUpdate("/accounts/alice")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not validate", func() {
			orig := string(hs.Object().Extra)
			Expect(subject.Patch(txn, "/accounts/alice", hs, &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{}`)},
			})).To(Succeed())
			Expect(hs.Object().Extra).To(MatchJSON(orig))
		})

		It("should hash password", func() {
			Expect(subject.Patch(txn, "/accounts/alice", hs, &schema.Resource{
				Data: &schema.Object{ID: "alice", Extra: []byte(`{"password":"upd@ted"}`)},
			})).To(Succeed())
			Expect(hs.Object().Extra).To(MatchJSON(`{"password":"$argon2id$v=19$m=32,t=1,p=1$Iw$Ef1fxv/3Imo"}`))
		})

		It("should retain account as writer", func() {
			Expect(subject.Patch(txn, "/accounts/alice", hs, &schema.Resource{
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
