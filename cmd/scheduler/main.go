package main

import (
	"math/rand"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"github.com/alice890308/scheduling-framwork/pkg/plugin"
	"k8s.io/component-base/logs"
	"os"
	"time"
	"fmt"
)

func main () {
	rand.Seed(time.Now().UTC().UnixNano())

	command := app.NewSchedulerCommand(
		app.WithPlugin(plugins.Name, plugins.New),
	)

	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}