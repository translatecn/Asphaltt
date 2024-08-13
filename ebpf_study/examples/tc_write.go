package examples

import (
	"ebpf_study/bpf"
	"errors"
	"fmt"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	"log"
	"time"
	"unsafe"
)

type Data struct {
	Pid  uint64
	Name [255]byte
}

func Load2() {
	err := rlimit.RemoveMemlock()
	if err != nil {
		log.Fatalln(err)
	}
	tcObj := bpf.TC_WRITEObjects{}
	err = bpf.LoadTC_WRITEObjects(&tcObj, nil)
	if err != nil {
		log.Fatalln(err)
	}
	// 不会阻塞
	tracepoint, err := link.Tracepoint("syscalls", "sys_enter_write", tcObj.HandleTp, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer tracepoint.Close()

	//reader, err := perf.NewReader(tcObj.MyBpfMap, os.Getpagesize())
	reader, err := ringbuf.NewReader(tcObj.LogMap)
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()

	go func() {
		for {
			record, err := reader.Read()
			if err != nil {
				if errors.Is(err, perf.ErrClosed) {
					log.Println("Receiver signal, exiting...")
					return
				}
				log.Println("reading from reader:", err)
				continue
			}
			pointer := (*Data)(unsafe.Pointer(&record.RawSample[0]))
			log.Println("Record:", pointer.Pid, string(pointer.Name[:]))
		}
	}()

	fmt.Println(tracepoint)
	time.Sleep(time.Minute)
}
