module github.com/cloudfoundry-incubator/disaster-recovery-acceptance-tests

go 1.16

require (
	code.cloudfoundry.org/cfhttp/v2 v2.0.0 // indirect
	code.cloudfoundry.org/routing-api v0.0.0-20200423182059-8cf9b4589b27
	code.cloudfoundry.org/trace-logger v0.0.0-20170119230301-107ef08a939d // indirect
	github.com/bmizerany/pat v0.0.0-20170815010413-6226ea591a40 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.12.0
	github.com/tedsuo/rata v1.0.0 // indirect
	github.com/vito/go-sse v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20210503173754-0981d6026fa6 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace credhub-test-app => ./fixtures/credhub-test-app
