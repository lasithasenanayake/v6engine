package main

import (
	"duov6.com/cebadapter"
	"duov6.com/objectstore/endpoints"
	"duov6.com/objectstore/unittesting"
	"fmt"
	"github.com/fatih/color"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var isUnitTestMode bool = false

	if isUnitTestMode {
		unittesting.Start()
	} else {
		splash()
		initialize()
	}
}

func initialize() {

	cebadapter.Attach("ObjectStore", func(s bool) {
		cebadapter.GetLatestGlobalConfig("StoreConfig", func(data []interface{}) {
			fmt.Println()
			fmt.Println(data)
			fmt.Println()
			color.Yellow("Store Configuration Successfully Loaded...")
			agent := cebadapter.GetAgent()

			agent.Client.OnEvent("globalConfigChanged.StoreConfig", func(from string, name string, data map[string]interface{}, resources map[string]interface{}) {
				cebadapter.GetLatestGlobalConfig("StoreConfig", func(data []interface{}) {
					color.Yellow("Store Configuration Successfully Updated...")
				})
			})
		})
		color.Yellow("Successfully registered in CEB")
	})

	httpServer := endpoints.HTTPService{}
	go httpServer.Start()

	bulkService := endpoints.BulkTransferService{}
	go bulkService.Start()

	forever := make(chan bool)
	<-forever
}

func splash() {
	color.Green("")
	color.Green("")
	color.Green("                                                 ~~")
	color.Green("    ____             _____ __                  | ][ |")
	color.Green("   / __ \\__  ______ / ___// /_____  ________     ~~")
	color.Green("  / / / / / / / __ \\__ \\/ __/ __ \\/ ___/ _ \\")
	color.Green(" / /_/ / /_/ / /_/ /__/ / /_/ /_/ / /  /  __/")
	color.Green("/_____/\\__,_/\\____/____/\\__/\\____/_/   \\___/ ")
	color.Green("")
	color.Green("")
	color.Green("")
}
