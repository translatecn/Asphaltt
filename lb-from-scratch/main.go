package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"net"

	"github.com/dropbox/goebpf"
	"github.com/spf13/cobra"
)

//go:embed ebpf_prog/xdp_lb.elf
var xdpProg []byte

type config struct {
	ifname  string
	forward string
}

var cfg config

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringVarP(&cfg.ifname, "ifname", "i", "", "interface to run the xdp lb")
	fs.StringVarP(&cfg.forward, "fwd", "o", "", "interface to send out the lb packet")
}

var rootCmd = cobra.Command{
	Use: "xdp-lb",
	Run: func(cmd *cobra.Command, args []string) {
		xdpLb()
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}

func xdpLb() {
	ifi, err := net.InterfaceByName(cfg.ifname)
	if err != nil {
		fmt.Println(cfg.ifname, "is not an interface")
		return
	}

	if err := loadXDP(ifi.Index); err != nil {
		fmt.Println("failed to setup xdp to", cfg.ifname, "err:", err)
	} else {
		fmt.Println("setup xdp to", cfg.ifname, "successfully")
	}
}

func loadXDP(ifindex int) error {
	bpf := goebpf.NewDefaultEbpfSystem()
	if err := bpf.Load(bytes.NewReader(xdpProg)); err != nil {
		return fmt.Errorf("failed to load xdp program, err: %w", err)
	}

	xdpRedirect, err := loadXDPProg(bpf)
	if err != nil {
		return fmt.Errorf("failed to load xdp program, err: %w", err)
	}

	if err := attachXDPProg(xdpRedirect, cfg.ifname); err != nil {
		return fmt.Errorf("failed to attach bpf prog to interface(%s), err: %w", cfg.ifname, err)
	}

	return nil
}

func loadXDPProg(bpf goebpf.System) (goebpf.Program, error) {
	xdpRedirect := bpf.GetProgramByName("xdp_load_balancer")
	if xdpRedirect == nil {
		return nil, fmt.Errorf("bpf prog(xdp_load_balancer) not found")
	}

	if err := xdpRedirect.Load(); err != nil {
		return nil, fmt.Errorf("failed to load bpf prog(xdp_load_balancer), err: %w", err)
	}

	return xdpRedirect, nil
}

func attachXDPProg(prog goebpf.Program, ifname string) error {
	param := &goebpf.XdpAttachParams{Interface: ifname, Mode: goebpf.XdpAttachModeDrv}
	err := prog.Attach(param)
	if err == nil {
		return nil
	}

	param.Mode = goebpf.XdpAttachModeSkb
	err = prog.Attach(param)
	if err != nil {
		return fmt.Errorf("failed to attach bpf prog to interface(%s), err: %w", ifname, err)
	}

	return nil
}
