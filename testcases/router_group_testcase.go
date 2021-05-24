package testcases

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	routing_api "code.cloudfoundry.org/routing-api"

	"code.cloudfoundry.org/routing-api/models"
	. "github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests/runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type none struct{}
type status bool

var cleanupRouterGroups func() error

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

func (tc *CfRouterGroupTestCase) CheckDeployment(config Config) {
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
	time.Sleep(15 * time.Second)
	tc.routingAPIClient = routing_api.NewClient(config.CloudFoundryConfig.APIURL, true)
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
		Guid:            "RandomTestGUID" + "_" + RandomStringNumber(),
		Name:            "RandomTestName" + "_" + RandomStringNumber(),
		Type:            "tcp",
		ReservablePorts: "1024-2047",
	}

	cleanupRouterGroups = func() error { return tc.routingAPIClient.DeleteRouterGroup(routerGroupEntry) }
	_, err := routerGroupRequestWithRetry(func() (models.RouterGroups, error) {
		return nil, tc.routingAPIClient.CreateRouterGroup(routerGroupEntry)
	}, 5)
	Expect(err).NotTo(HaveOccurred())
}

func (tc *CfRouterGroupTestCase) EnsureAfterSelectiveRestore(config Config) {}

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
	time.Sleep(15 * time.Second)
	routerGroupsPostRestore, err := tc.readRouterGroups(token)
	Expect(err).NotTo(HaveOccurred())
	Expect(routerGroupsPostRestore).To(ConsistOf(tc.routerGroupsPreBackup))
}

// Cleanup is called at the end to remove the test artifacts left behind.
func (tc *CfRouterGroupTestCase) Cleanup(config Config) {
	cleanupRouterGroups()
}

func loginAndGetToken(config Config) string {
	RunCommandSuccessfully("cf login --skip-ssl-validation -a", config.CloudFoundryConfig.APIURL, "-u", config.CloudFoundryConfig.AdminUsername, "-p", config.CloudFoundryConfig.AdminPassword)
	return refreshToken()
}

func refreshToken() string {
	token := string(RunCommandSuccessfullySilently("cf oauth-token").Out.Contents()[:])
	token = strings.Split(token, " ")[1]
	token = strings.Trim(token, "\r\n\t ")
	return token
}

func (tc *CfRouterGroupTestCase) readRouterGroups(token string) (models.RouterGroups, error) {
	tc.routingAPIClient.SetToken(token)
	response, err := routerGroupRequestWithRetry(func() (models.RouterGroups, error) {
		return tc.routingAPIClient.RouterGroups()
	}, 5)
	return response, err
}

func routerGroupRequestWithRetry(request func() (models.RouterGroups, error), retries int) (response models.RouterGroups, err error) {
	response, err = request()
	for attempt := 0; attempt < retries; attempt += 1 {
		if err != nil {
			switch err.(type) {
			case *url.Error:
				time.Sleep(time.Duration(attempt*attempt) * time.Second)
				fmt.Printf("--- RouterGroup request retry attempt:%v ---\n\n", attempt)
				response, err = request()
			default:
				break
			}
		} else {
			break
		}
	}
	return response, err
}
