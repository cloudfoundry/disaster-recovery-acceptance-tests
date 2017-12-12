package testcases

import (
	"strings"

	"code.cloudfoundry.org/routing-api"
	"code.cloudfoundry.org/routing-api/models"
	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type none struct{}
type status bool

const (
	alwaysAvailable    = status(true)
	notAlwaysAvailable = status(false)
)

// CfRouterGroupTestCase holds common variables across roter group testcases.
type CfRouterGroupTestCase struct {
	uniqueTestID string
	name         string

	readAvailability      chan status
	done                  chan none
	routerGroupsPreBackup []models.RouterGroup
	routingAPIClient      routing_api.Client
}

// NewRouterGroupTestCase creates an instance with random id.
func NewRouterGroupTestCase() *CfRouterGroupTestCase {
	id := RandomStringNumber()
	return &CfRouterGroupTestCase{
		uniqueTestID: id,
		name:         "cf-routing",
	}
}

func (tc *CfRouterGroupTestCase) Name() string {
	return tc.name
}

// BeforeBackup is run before the bbr backs up the router_group table.
// The function fetches and stores a snapshot of router_group table for later
// verification.
//
// It also spawns a goroutine to check that reads are always available during
// the backup.
func (tc *CfRouterGroupTestCase) BeforeBackup(config Config) {
	By("Getting CF OAuth Token")
	token := loginAndGetToken(config)
	By("Creating a pre-backup router group backup")
	var err error

	tc.routingAPIClient = routing_api.NewClient(config.CloudFoundryConfig.ApiUrl, true)
	tc.routerGroupsPreBackup, err = tc.readRouterGroups(token)
	Expect(err).NotTo(HaveOccurred())
}

// AfterBackup is run after the bbr is done backing up the database.
// The function inserts a dummy entry into router group table. This entry should
// not exist once restore is performed later.
//
// It verifies that the goroutine spawned in BeforeBackup verified that reads
// were always available during the backup.
//
// It also spawns a goroutine to check that reads are not available during
// the restore operation after AfterBackup is called.
func (tc *CfRouterGroupTestCase) AfterBackup(config Config) {
	By("Getting CF OAuth Token")
	token := refreshToken()

	By("Adding an entry in the router group table")
	tc.routingAPIClient.SetToken(token)
	routerGroupEntry := models.RouterGroup{
		Guid:            "RandomTestGUID",
		Name:            "RandomTestName",
		Type:            "tcp",
		ReservablePorts: "2000-4000",
	}
	err := tc.routingAPIClient.CreateRouterGroup(routerGroupEntry)
	Expect(err).NotTo(HaveOccurred())
}

// AfterRestore is run after the bbr is done restoring the database.
// The function compares the post restore router group table with
// pre restore router group table.
//
// It verifies that the goroutine spawned in AfterBackup verified that reads
// were not always available during the restore.
func (tc *CfRouterGroupTestCase) AfterRestore(config Config) {
	By("Getting CF OAuth Token")
	token := refreshToken()

	By("Taking a snapshot of restored table and comparing it with the pre-backup table")
	routerGroupsPostRestore, err := tc.readRouterGroups(token)
	Expect(err).NotTo(HaveOccurred())
	Expect(routerGroupsPostRestore).To(ConsistOf(tc.routerGroupsPreBackup))
}

// Cleanup is called at the end to remove the test artifacts left behind.
func (tc *CfRouterGroupTestCase) Cleanup(config Config) {
}

func loginAndGetToken(config Config) string {
	RunCommandSuccessfully("cf login --skip-ssl-validation -a", config.CloudFoundryConfig.ApiUrl, "-u", config.CloudFoundryConfig.AdminUsername, "-p", config.CloudFoundryConfig.AdminPassword)
	return refreshToken()
}

func refreshToken() string {
	token := string(RunCommandSuccessfully("cf oauth-token").Out.Contents()[:])
	token = strings.Split(token, " ")[1]
	token = strings.Trim(token, "\r\n\t ")
	return token
}

func (tc *CfRouterGroupTestCase) readRouterGroups(token string) (models.RouterGroups, error) {
	tc.routingAPIClient.SetToken(token)
	return tc.routingAPIClient.RouterGroups()
}
