package main

import (
	"context"
	"fmt"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type In struct {
	Name string
}

type Out struct {
	Name   string
	Result string
}

var HostPort = "127.0.0.1:7933"
var Domain = "domain-harish"
var TaskListName = "SimpleWorker"
var CadenceService = "cadence-frontend"
var ClientName = "harish-client"

func buildLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	var err error
	logger, err := config.Build()
	if err != nil {
		panic("Failed to setup logger")
	}

	return logger
}

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

	cl := client.NewClient(
		service,
		Domain,
		&client.Options{
			FeatureFlags: client.FeatureFlags{
				WorkflowExecutionAlreadyCompletedErrorEnabled: true,
			},
		})

	for i := 0; i < 100; i++ {
		go func(_i int) {
			for {
				run, err := cl.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
					TaskList:                        TaskListName,
					ExecutionStartToCloseTimeout:    10 * time.Second,
					DecisionTaskStartToCloseTimeout: time.Minute,
				}, "foo", &In{Name: "Harish"})
				if err != nil {
					panic(err)
				}
				var out Out
				err = run.Get(context.Background(), &out)
				if err != nil {
					panic(err)
				}
				fmt.Println(out)
				time.Sleep(100 * time.Millisecond)
			}
		}(i)

	}

	time.Sleep(100 * time.Hour)
}
