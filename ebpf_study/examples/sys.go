package examples

import (
	"ebpf_study/bpf"
	"errors"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/btf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/ringbuf"
	"log"
	"time"
	"unsafe"
)

type Proc struct {
	Pid  uint64
	Name [255]byte
}

func Load3() {
	sysObj := bpf.SYSObjects{}
	err := bpf.LoadSYSObjects(&sysObj, nil)
	if err != nil {
		log.Fatalln(err)
	}
	//tracepoint, err := link.Tracepoint("syscalls", "sys_enter_execve", sysObj.HandleTp, nil)
	tracepoint, err := link.Tracepoint("syscalls", "sys_exit_execve", sysObj.HandleTp, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer tracepoint.Close()

	// 创建 reader 读取 内核map

	reader, err := ringbuf.NewReader(sysObj.ProcMap)
	if err != nil {
		log.Fatalln(err)
	}

	defer reader.Close()

	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, perf.ErrClosed) {
				log.Println("received signal, exiting...")
				return
			}
			log.Println("received signal, exiting...")

			continue
		}
		if len(record.RawSample) > 0 {
			data := (*Proc)(unsafe.Pointer(&record.RawSample[0]))
			log.Println(data)
		}
	}

	time.Sleep(time.Minute)
}

func Multi() {
	rc := map[string]interface{}{
		"XDPACL_DEBUG":                   1,
		"XDPACL_BITMAP_ARRAY_SIZE_LIMIT": uint32(1 << 12),
	}
	//btfSpec, _ := btf.LoadSpec(flags.KernelBTF)
	btfSpec, _ := btf.LoadKernelSpec()

	bpfSpec, _ := loadSpec(8)
	if err := bpfSpec.RewriteConstants(rc); err != nil {
	}
	var opts ebpf.CollectionOptions
	opts.Programs.KernelTypes = btfSpec
	bpfSpec.LoadAndAssign(new(bpf.XDPACL8Objects), &opts)
}
