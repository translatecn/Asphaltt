package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type event memory ../c/memory/memory.c
func main() {
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)

	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	// Load pre-compiled programs and maps into the kernel.
	objs := memoryObjects{}
	if err := loadMemoryObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %s", err)
	}
	defer objs.Close()

	kp, err := link.Tracepoint("kmem", "kmalloc", objs.TraceKmalloc, nil)
	if err != nil {
		log.Fatalf("opening tracepoint: %s", err)
	}
	defer kp.Close()

	reader, err := ringbuf.NewReader(objs.Events)
	if err != nil {
		log.Fatalf("creating perf event reader: %s", err)
	}

	go func() {
		// Wait for a signal and close the perf reader,
		// which will interrupt rd.Read() and make the program exit.
		<-stopper
		log.Println("Received signal, exiting program........................................................")

		// 关闭reader
		err := reader.Close()
		if err != nil {
			log.Fatalf("closing reader: %s", err)
		}

		if err := kp.Close(); err != nil {
			log.Fatalf("closing perf event reader: %s", err)
		}
	}()

	var event memoryEvent
	for true {
		record, err := reader.Read()
		if err != nil {
			return
		}
		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("parsing perf event: %s", err)
			continue
		}
		println("pid:" + strconv.Itoa(int(event.Pid)) + " size:" + strconv.Itoa(int(event.Size)))
	}
}
