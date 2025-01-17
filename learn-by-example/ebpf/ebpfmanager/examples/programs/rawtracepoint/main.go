package main

import (
	"github.com/sirupsen/logrus"

	manager "github.com/Asphaltt/go-nfnetlink-example/ebpf/ebpfmanager"
)

var m = &manager.Manager{
	Probes: []*manager.Probe{
		{
			Section:      "raw_tracepoint/sys_enter",
			EbpfFuncName: "raw_tracepoint_sys_enter",
		},
	},
}

func main() {
	// Initialize the manager
	if err := m.Init(recoverAssets()); err != nil {
		logrus.Fatal(err)
	}

	// Start the manager
	if err := m.Start(); err != nil {
		logrus.Fatal(err)
	}

	logrus.Println("successfully started, head over to /sys/kernel/debug/tracing/trace_pipe")

	// Create a folder to trigger the probes
	if err := trigger(); err != nil {
		logrus.Error(err)
	}

	// Close the manager
	if err := m.Stop(manager.CleanAll); err != nil {
		logrus.Fatal(err)
	}
}
