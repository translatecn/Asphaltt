package main

import (
	"bytes"
	"log"
	"net"

	"github.com/dropbox/goebpf"
	"github.com/spf13/cobra"
)

func init() {
	var dummyCmd = cobra.Command{
		Use: "dummy",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfg.ifnames) != 1 {
				log.Println("a dummy interface should be provided")
				return
			}

			dummyXDP()
		},
	}

	rootCmd.AddCommand(&dummyCmd)
}

func dummyXDP() {
	ifname := cfg.ifnames[0]
	_, err := net.InterfaceByName(ifname)
	if err != nil {
		log.Printf("%s is not an interface\n", ifname)
		return
	}

	bpf := goebpf.NewDefaultEbpfSystem()
	if err := bpf.Load(bytes.NewReader(xdpProg)); err != nil {
		log.Println("failed to load bpf program, err:", err)
		return
	}

	dummy := bpf.GetProgramByName("xdp_pass")
	if dummy == nil {
		log.Println("bpf prog(xdp_pass) not found")
		return
	}

	if err := dummy.Load(); err != nil {
		log.Println("failed to load bpf prog(xdp_pass), err:", err)
		return
	}

	if err := attachXDPProg(dummy, ifname); err != nil {
		log.Println("failed to attach bpf prog(xdp_pass) to", ifname, "err:", err)
		return
	}
}
