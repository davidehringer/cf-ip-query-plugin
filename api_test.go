package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/cli/plugin/fakes"
   "io/ioutil"
)

var _ = Describe("API", func(){
	Describe("StartedApps", func(){
		var cliConnection *fakes.FakeCliConnection

		BeforeEach(func() {
			cliConnection = &fakes.FakeCliConnection{}
		})

		It("Lists app starts apps", func(){
         getAppsJsonPage, _ := ioutil.ReadFile("fixtures/apps.json")
			cliConnection.CliCommandWithoutTerminalOutputReturns([]string{string(getAppsJsonPage)}, nil)

			api := NewCliCcApi(cliConnection)
			apps, _ := api.StartedApps()

			Expect(len(apps)).To(Equal(45))
		})
	})

   Describe("AppStats", func(){
      var cliConnection *fakes.FakeCliConnection

      BeforeEach(func() {
         cliConnection = &fakes.FakeCliConnection{}
      })

      It("Gets basic app stats", func(){
         appStatsJson, _ := ioutil.ReadFile("fixtures/app-stats.json")
         cliConnection.CliCommandWithoutTerminalOutputReturns([]string{string(appStatsJson)}, nil)

         api := NewCliCcApi(cliConnection)
         app, _ := api.AppStats("aca81d95-e489-4789-87c8-c2444f50f1e3")

         Expect(len(app.AppInstanceStats)).To(Equal(4))
      })

      It("When all instances are down", func(){
         appStatsJson, _ := ioutil.ReadFile("fixtures/app-stats-down.json")
         cliConnection.CliCommandWithoutTerminalOutputReturns([]string{string(appStatsJson)}, nil)

         api := NewCliCcApi(cliConnection)
         app, _ := api.AppStats("aca81d95-e489-4789-87c8-c2444f50f1e3")

         Expect(len(app.AppInstanceStats)).To(Equal(0))
      })
   })
})