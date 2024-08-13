all: build setup mapid

clean:
	rm -f iptables-bpf.elf
	rm -f iptables-bpf
	iptables -D OUTPUT -m bpf --object-pinned /sys/fs/bpf/iptbpf -j DROP
	rm -f /sys/fs/bpf/iptbpf

rebuild: clean build

build:
	clang -I../headers -O2 -g -target bpf -c iptables-bpf.c -o iptables-bpf.elf
	go build -v -o iptables-bpf main.go

setup: build
	bpftool prog load iptables-bpf.elf /sys/fs/bpf/iptbpf
	iptables -I OUTPUT -m bpf --object-pinned /sys/fs/bpf/iptbpf -j DROP

mapid:
	@bpftool map list | grep filter_daddrs | awk -F: '{print $$1}'
	@bpftool map list | grep filter_daddrs | awk -F: '{print $$1}'|xargs -I F ./iptables-bpf -m F -d 114.114.114.114

