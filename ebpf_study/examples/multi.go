package examples

import (
	"ebpf_study/bpf"
	"fmt"
	"github.com/cilium/ebpf"
)

func loadSpec(bitmapArraySize int) (*ebpf.CollectionSpec, error) {
	var err error
	var bpfSpec *ebpf.CollectionSpec
	switch n := bitmapArraySize; n {
	case 8:
		bpfSpec, err = bpf.LoadXDPACL8()
	case 16:
		bpfSpec, err = bpf.LoadXDPACL16()
	case 32:
		bpfSpec, err = bpf.LoadXDPACL32()
	case 64:
		bpfSpec, err = bpf.LoadXDPACL64()
	case 128:
		bpfSpec, err = bpf.LoadXDPACL128()
	case 160:
		bpfSpec, err = bpf.LoadXDPACL160()
	case 256:
		bpfSpec, err = bpf.LoadXDPACL256()
	default:
		err = fmt.Errorf("no bpf spec for %d", n)
	}
	return bpfSpec, err
}
