package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cilium/ebpf"
)

func must[T any](x T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return x
}

// 使用 Go 代码来查看采用了 CO-RE 编译的 eBPF 汇编

func main() {
	//if len(os.Args) != 2 {
	//	log.Fatal("miss elf file")
	//}
	//elfFile := os.Args[1]
	elfFile := `/Users/acejilam/Desktop/todo/Asphaltt/learn-by-example/tools/asm/tcpconn_bpfeb.o`
	fd := must(os.Open(elfFile))
	defer fd.Close()

	spec := must(ebpf.LoadCollectionSpecFromReader(fd))
	for _, prog := range spec.Programs {
		fmt.Printf("%v\n", prog.Instructions)
	}
}
