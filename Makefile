LLVM_PATH ?= /usr/bin/

CLANG ?= $(LLVM_PATH)/clang
STRIP ?= $(LLVM_PATH)/llvm-strip
CFLAGS := -O2 -g -Wall -Werror $(CFLAGS)
GOOS := linux
GOLDFLAGS := -s -w

bpf/bpf_bpfel.go: export BPF_STRIP := $(STRIP)
bpf/bpf_bpfel.go: export BPF_CLANG := $(CLANG)
bpf/bpf_bpfel.go: export BPF_CFLAGS := $(CFLAGS)
bpf/bpf_bpfel.go: bpf/xdp.c
	go generate ./...

.PHONY: generate
generate: bpf/bpf_bpfel.go

sysperf: export GOOS := $(GOOS)
sysperf: generate
	go build -ldflags "$(GOLDFLAGS)"

clean:
	@rm -f bpf/bpf_* sysperf

.DEFAULT_GOAL := sysperf
