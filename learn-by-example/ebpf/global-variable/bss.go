package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/btf"
	"github.com/iovisor/gobpf/pkg/bpffs"
)

// https://mp.weixin.qq.com/s?__biz=MjM5MTQxNTk5MA==&mid=2247484091&idx=1&sn=eee926d16e80fa944f5daa2b972dca81&scene=21#wechat_redirect

//对比					全局变量						常量
//声明:					__be32 filter_addr;			volatile const __be32 filter_addr;
//使用:					直接使用						直接使用
//更新:					运行时						加载时
//bpf map 名字:		.bss 或者 .data					.rodata

const bssMapName = "gvar_bss"

func loadBssMap(spec *ebpf.MapSpec) *ebpf.Map {
	mapPinPath := filepath.Join(bpffs.BPFFSPath, bssMapName)
	if m, err := ebpf.LoadPinnedMap(mapPinPath, nil); err == nil {
		return m
	}

	spec.Name = bssMapName
	spec.Pinning = ebpf.PinByName
	m, err := ebpf.NewMapWithOptions(spec, ebpf.MapOptions{
		PinPath: bpffs.BPFFSPath,
	})
	if err != nil {
		log.Fatalf("Failed to new bpf map %s: %v", bssMapName, err)
	}

	return m
}

func genBssStruct() {
	bpfSpec, err := loadTcpconn()
	if err != nil {
		log.Fatalf("Failed to load bpf spec: %v", err)
	}

	m, ok := bpfSpec.Maps[".bss"] // 常亮 .rodata 加载时赋值
	if !ok {
		log.Fatalf(".bss map not found")
	}

	fmt.Printf(".bss map spec: %v\n", m)

	var gof btf.GoFormatter
	out, err := gof.TypeDeclaration("bssValue", m.Value)
	if err != nil {
		log.Fatalf("Failed to generate Go struct for .bss value")
	}
	// output:
	// root@leonhwang-svr ~/P/ebpf-globalvar# ./ebpf-globalvar -gen
	// .bss map spec: Array(keySize=4, valueSize=6, maxEntries=1, flags=0)
	// .bss map value: type bssValue struct { filter_daddr uint32; filter_dport uint16; }
	fmt.Println(".bss map value:", out)
}

// get and update from genBssStruct()
type bssValue struct {
	Daddr [4]byte
	Dport uint16
}
