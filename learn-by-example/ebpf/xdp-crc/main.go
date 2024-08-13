// Copyright 2024 Leon Hwang.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	flag "github.com/spf13/pflag"
	"github.com/vishvananda/netlink"
)

//go:generate bpf2go -cc clang xdp ./xdp.c -- -D__TARGET_ARCH_x86 -I../headers -Wall

func main() {
	var device string
	flag.StringVarP(&device, "device", "d", "lo", "device to attach XDP program")
	flag.Parse()

	ifi, _ := netlink.LinkByName(device)

	rlimit.RemoveMemlock()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var obj xdpObjects
	loadXdpObjects(&obj, nil)
	defer obj.Close()

	xdp, _ := link.AttachXDP(link.XDPOptions{
		Program:   obj.Crc,
		Interface: ifi.Attrs().Index,
	})
	defer xdp.Close()

	log.Printf("Attached xdp to %s", device)

	<-ctx.Done()
}
