package main

import (
	"bytes"
	"fmt"
	"net"

	"github.com/dropbox/goebpf"
	"github.com/spf13/cobra"
)

func init() {
	var dummyCmd = cobra.Command{
		Use: "dummy",
		Run: func(cmd *cobra.Command, args []string) {
			dummyXDP()
		},
	}

	rootCmd.AddCommand(&dummyCmd)
}

func dummyXDP() {
	ifname := cfg.ifname
	_, err := net.InterfaceByName(ifname)
	if err != nil {
		fmt.Printf("%s is not an interface\n", ifname)
		return
	}

	bpf := goebpf.NewDefaultEbpfSystem()
	if err := bpf.Load(bytes.NewReader(xdpProg)); err != nil {
		fmt.Println("failed to load bpf program, err:", err)
		return
	}

	dummy := bpf.GetProgramByName("xdp_pass")
	if dummy == nil {
		fmt.Println("bpf prog(xdp_pass) not found")
		return
	}

	if err := dummy.Load(); err != nil {
		fmt.Println("failed to load bpf prog(xdp_pass), err:", err)
		return
	}

	if err := attachXDPProg(dummy, ifname); err != nil {
		fmt.Println("failed to attach bpf prog(xdp_pass) to", ifname, "err:", err)
		return
	}
}
