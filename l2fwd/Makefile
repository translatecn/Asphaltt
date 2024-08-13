GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
CLANG := clang
CLANG_INCLUDE := -I./ebpf/headers

GO_SOURCE := ./*.go
GO_BINARY := l2fwd

EBPF_SOURCE := ebpf/xdp.c
EBPF_BINARY := ebpf/xdp.elf


.PHONY: all debug rebuild build_ebpf build_go clean

all: build_ebpf build_go

debug:
	$(CLANG) $(CLANG_INCLUDE) -O2 -g -target bpf -c $(EBPF_SOURCE)  -o $(EBPF_BINARY) -DDEBUG
	$(GOBUILD) -v -o $(GO_BINARY) $(GO_SOURCE)

rebuild: clean all

build_ebpf: $(EBPF_BINARY)

build_go: $(GO_BINARY)

clean:
	$(GOCLEAN)
	rm -f $(GO_BINARY)
	rm -f $(EBPF_BINARY)

$(EBPF_BINARY): $(EBPF_SOURCE)
	$(CLANG) $(CLANG_INCLUDE) -O2 -g -target bpf -c $^  -o $@
	rm -f $(GO_BINARY)

$(GO_BINARY): $(GO_SOURCE)
	$(GOBUILD) -v -o $@ $^
