package main

import (
	"cadence-example/worker/echo"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"time"
)

var HostPort = "127.0.0.1:7933"
var Domain = "domain-harish"
var TaskListName = "SimpleWorker"
var CadenceService = "cadence-frontend"
var ClientName = "harish-client"

func main() {

	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(ClientName))
	if err != nil {
		panic("Failed to setup tchannel")
	}
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: ClientName,
		Outbounds: yarpc.Outbounds{
			CadenceService: {Unary: ch.NewSingleOutbound(HostPort)},
		},
	})
	if err := dispatcher.Start(); err != nil {
		panic("Failed to start dispatcher")
	}

	service := workflowserviceclient.New(dispatcher.ClientConfig(CadenceService))
	workerObj := worker.New(service, Domain, TaskListName,
		worker.Options{
			StickyScheduleToStartTimeout: 10,
		},
	)
	workerObj.RegisterWorkflowWithOptions(echo.Workflow, workflow.RegisterOptions{Name: "foo"})
	workerObj.RegisterActivity(echo.Activity)

	err = workerObj.Start()
	if err != nil {
		panic(err)
	}

	err = workerObj.Run()
	if err != nil {
		panic(err)
	}

	time.Sleep(24 * time.Hour)
}
