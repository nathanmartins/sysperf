package custom_ebpf

// This program demonstrates attaching an eBPF program to a kernel symbol and
// using percpu map to collect data. The eBPF program will be attached to the
// start of the sys_execve kernel function and prints out the number of called
// times on each cpu every second.

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

// $BPF_CLANG and $BPF_CFLAGS are set by the Makefile.
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf kprobe_percpu.c -- -I./headers

const mapKey uint32 = 0

type CPUInfo struct {
	FunctionName string
	Times        int
	Id           int
}

func RunKprobe(ctx context.Context) error {

	// Name of the kernel function to trace.
	fn := "sys_execve"

	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	// Load pre-compiled programs and maps into the kernel.
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	// Open a Kprobe at the entry point of the kernel function and attach the
	// pre-compiled program. Each time the kernel function enters, the program
	// will increment the execution counter by 1. The read loop below polls this
	// map value once per second.
	kp, err := link.Kprobe(fn, objs.KprobeExecve)
	if err != nil {
		log.Fatalf("opening kprobe: %s", err)
	}
	defer kp.Close()

	// Read loop reporting the total amount of times the kernel
	// function was entered, once per second.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var allCpuValue []uint64
			if err = objs.KprobeMap.Lookup(mapKey, &allCpuValue); err != nil {
				return fmt.Errorf("reading map: %v", err)
			}
			for cpuId, cpuValue := range allCpuValue {
				//log.Printf("%s called %d times on CPU%v\n", fn, cpuValue, cpuId)

				i := CPUInfo{
					FunctionName: fn,
					Times:        int(cpuValue),
					Id:           cpuId,
				}

				var output []byte
				output, err = json.Marshal(i)
				if err != nil {
					return err
				}

				fmt.Println(string(output))
			}
		case <-ctx.Done():
			return nil
		}
	}

	return err
}
