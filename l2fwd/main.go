package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/dropbox/goebpf"
	"github.com/spf13/cobra"
)

//go:embed ebpf/xdp.elf
var xdpProg []byte

type config struct {
	ifnames []string
}

var cfg config

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringArrayVarP(&cfg.ifnames, "ifindex", "I", nil, "ifindex array to interchange layer 2 packets")
}

var rootCmd = cobra.Command{
	Use: "l2fwd",
	Run: func(cmd *cobra.Command, args []string) {
		ifnames := cfg.ifnames
		if len(ifnames) > 0 && len(ifnames) != 2 {
			ifnames = strings.Split(ifnames[0], ",")
		}
		if len(ifnames) != 2 {
			fmt.Printf("at least 2 ifindex are required, got %d\n", len(cfg.ifnames))
			return
		}

		cfg.ifnames = ifnames
		l2fwd()
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}

func l2fwd() {
	var ifindex []int
	for _, name := range cfg.ifnames {
		ifi, err := net.InterfaceByName(name)
		if err != nil {
			log.Println(name, "is not an interface")
			return
		}

		log.Printf("ifindex(%d) is for ifname(%s)\n", ifi.Index, name)

		ifindex = append(ifindex, ifi.Index)
	}

	if err := loadXDP(ifindex); err != nil {
		log.Printf("failed to load xdp program, err: %v", err)
	}
}

func loadXDP(ifindex []int) error {
	bpf := goebpf.NewDefaultEbpfSystem()
	if err := bpf.Load(bytes.NewReader(xdpProg)); err != nil {
		return fmt.Errorf("failed to load xdp program, err: %w", err)
	}

	if err := prepareXDP(bpf, ifindex); err != nil {
		return fmt.Errorf("failed to prepare xdp program, err: %w", err)
	}

	xdpRedirect, err := loadXDPProg(bpf)
	if err != nil {
		return fmt.Errorf("failed to load xdp program, err: %w", err)
	}

	if err := attachXDPProg(xdpRedirect, cfg.ifnames[0]); err != nil {
		return fmt.Errorf("failed to attach bpf prog to interface(%s), err: %w", cfg.ifnames[0], err)
	}

	if err := attachXDPProg(xdpRedirect, cfg.ifnames[1]); err != nil {
		return fmt.Errorf("failed to attach bpf prog to interface(%s), err: %w", cfg.ifnames[1], err)
	}

	return nil
}

func prepareXDP(bpf goebpf.System, ifindex []int) error {
	if len(ifindex) != 2 {
		return fmt.Errorf("number of interface should be 2")
	}

	ingress, egress := ifindex[0], ifindex[1]

	mFdb := bpf.GetMapByName("m_fdb")
	if mFdb == nil {
		return fmt.Errorf("bpf map(m_fdb) not found")
	}

	if err := mapUpsert(mFdb, ingress, egress); err != nil {
		return fmt.Errorf("failed to update bpf map(m_fdb), err: %w", err)
	}

	mPorts := bpf.GetMapByName("m_ports")
	if mPorts == nil {
		return fmt.Errorf("bpf map(m_ports) not found")
	}

	if err := mPorts.Upsert(ingress, ingress); err != nil {
		fmt.Errorf("failed to update bpf map(m_ports), err: %w", err)
	}

	if err := mPorts.Upsert(egress, egress); err != nil {
		fmt.Errorf("failed to update bpf map(m_ports), err: %w", err)
	}

	return nil
}

func mapUpsert(m goebpf.Map, ingress, egress int) error {
	if err := m.Upsert(ingress, egress); err != nil {
		return err
	}
	return m.Upsert(egress, ingress)
}

func loadXDPProg(bpf goebpf.System) (goebpf.Program, error) {
	xdpRedirect := bpf.GetProgramByName("xdp_redirect")
	if xdpRedirect == nil {
		return nil, fmt.Errorf("bpf prog(xdp_redirect) not found")
	}

	if err := xdpRedirect.Load(); err != nil {
		return nil, fmt.Errorf("failed to load bpf prog(xdp_redirect), err: %w", err)
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
