package apps

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vito/cmdtest/matchers"

	. "../helpers"
)

var _ = Describe("Riak CS Service", func() {
	BeforeEach(func() {
		AppName = RandomName()

		Expect(Cf("push", AppName, "-m", "256M", "-p", sinatraPath, "-no-start")).To(Say("files\nOK"))
	})

	AfterEach(func() {
		Expect(Cf("delete", AppName, "-f")).To(Say("OK"))
	})

	It("Allows users to create, bind, write to, read from, unbind, and destroy the service instance", func() {
		ServiceName := "riak-cs"
		PlanName := "bucket"
		ServiceInstanceName := RandomName()

		Expect(Cf("create-service", ServiceName, PlanName, ServiceInstanceName)). To(Say("OK"))
		Expect(Cf("bind-service", AppName, ServiceInstanceName)).To(Say("OK"))
		Expect(Cf("start", AppName)).To(Say("App started"))

		uri := AppUri(AppName) + "/service/blobstore/" + ServiceInstanceName + "/mykey"
		delete_uri := AppUri(AppName) + "/service/blobstore/" + ServiceInstanceName

		Eventually(Curling("-d", "myvalue", uri), 10.0, 1.0).Should(Say("myvalue"))
		Eventually(Curling(uri), 10.0, 1.0).Should(Say("myvalue"))
		Eventually(Curling("-X", "DELETE", delete_uri), 10.0, 1.0).Should(Say(""))

		Expect(Cf("unbind-service", AppName, ServiceInstanceName)).To(Say("OK"))
		Expect(Cf("delete-service", "-f", ServiceInstanceName)). To(Say("OK"))
	})
})