LLVM_PATH ?= /usr/bin

CLANG ?= $(LLVM_PATH)/clang
STRIP ?= $(LLVM_PATH)/llvm-strip
CFLAGS := -O2 -g -Wall -Werror $(CFLAGS)
GOOS := linux
GOLDFLAGS := -s -w

custom_ebpf/bpf_bpfel.go: export BPF_STRIP := $(STRIP)
custom_ebpf/bpf_bpfel.go: export BPF_CLANG := $(CLANG)
custom_ebpf/bpf_bpfel.go: export BPF_CFLAGS := $(CFLAGS)
custom_ebpf/bpf_bpfel.go: custom_ebpf/kprobe_percpu.c
	go generate ./...

.PHONY: generate
generate: custom_ebpf/bpf_bpfel.go

sysperf: export GOOS := $(GOOS)
sysperf: generate
	go build -ldflags "$(GOLDFLAGS)"

clean:
	@rm -f custom_ebpf/bpf_* sysperf

.DEFAULT_GOAL := sysperf
