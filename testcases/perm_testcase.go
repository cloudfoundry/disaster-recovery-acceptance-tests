package testcases

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	"code.cloudfoundry.org/perm/pkg/perm"

	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type PermTestCase struct {
	uniqueTestID string
	roleName     string
	actor        perm.Actor
	permission   perm.Permission
	name         string
}

func NewPermTestCase() *PermTestCase {
	id := RandomStringNumber()
	actor := perm.Actor{
		ID:        id,
		Namespace: "actor_" + id,
	}
	permission := perm.Permission{
		Action:          "permission_" + id,
		ResourcePattern: "resource_" + id,
	}
	return &PermTestCase{
		uniqueTestID: id,
		roleName:     "test_role_" + id,
		actor:        actor,
		permission:   permission,
		name:         "perm",
	}
}

func (tc *PermTestCase) Name() string {
	return tc.name
}

func (tc *PermTestCase) BeforeBackup(config Config) {
	client := tc.permDial(config)
	defer client.Close()

	_, err := client.CreateRole(context.Background(), tc.roleName, tc.permission)
	Expect(err).NotTo(HaveOccurred())

	By("Creating an actor with a given role")
	err = client.AssignRole(context.Background(), tc.roleName, tc.actor)
	Expect(err).NotTo(HaveOccurred())

	By("Verifying that the actor has the given permission")
	ok, err := client.HasPermission(context.Background(), tc.actor, tc.permission.Action, tc.permission.ResourcePattern)
	Expect(err).NotTo(HaveOccurred())
	Expect(ok).To(BeTrue())
}

func (tc *PermTestCase) AfterBackup(config Config) {
	client := tc.permDial(config)
	defer client.Close()

	By("Deleting the actor")
	err := client.DeleteRole(context.Background(), tc.roleName)
	Expect(err).NotTo(HaveOccurred())

	By("Checking that the permission no longer exists")
	ok, err := client.HasPermission(context.Background(), tc.actor, tc.permission.Action, tc.permission.ResourcePattern)
	Expect(err).NotTo(HaveOccurred())
	Expect(ok).To(BeFalse())
}

func (tc *PermTestCase) AfterRestore(config Config) {
	client := tc.permDial(config)
	defer client.Close()

	By("Checking that the actor and its permission have been restored")
	ok, err := client.HasPermission(context.Background(), tc.actor, tc.permission.Action, tc.permission.ResourcePattern)
	Expect(err).NotTo(HaveOccurred())
	Expect(ok).To(BeTrue())
}

func (tc *PermTestCase) Cleanup(config Config) {
	client := tc.permDial(config)
	defer client.Close()

	err := client.DeleteRole(context.Background(), tc.roleName)
	Expect(err).NotTo(HaveOccurred())
}

func (tc *PermTestCase) permDial(config Config) *perm.Client {
	rootCAPool := x509.NewCertPool()

	ok := rootCAPool.AppendCertsFromPEM([]byte(config.PermCA))
	Expect(ok).To(BeTrue())

	client, err := perm.Dial(
		config.PermUrl,
		perm.WithTLSConfig(&tls.Config{RootCAs: rootCAPool}),
	)
	Expect(err).NotTo(HaveOccurred())

	return client
}
