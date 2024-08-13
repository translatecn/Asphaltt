CONFIG_STACK_VALIDATION=

ccflags-y  := -I.

obj-m += hello.o
obj-m += who-connect-me.o
obj-m += add-arp-records.o
obj-m += excited_virus.o
obj-m += check-tcp-syncookies.o
obj-m += custom-netlink.o
obj-m += kprobe_tcp_conn_request.o
obj-m += run-bpf-prog.o

all:
	mkdir /tmp/kernel-module-fun || echo "mkdir"
	mkdir ./ko || echo "mkdir"
	cp -rf ./* /tmp/kernel-module-fun
	make -C /lib/modules/$(shell uname -r)/build M=/tmp/kernel-module-fun modules
	cp /tmp/kernel-module-fun/*.ko ./ko

clean:
	make -C /lib/modules/$(shell uname -r)/build M=/tmp/kernel-module-fun clean
