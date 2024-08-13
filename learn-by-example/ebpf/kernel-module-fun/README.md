Kernel module fun
=================

## Motivation

I didn't know at all how kernel modules worked. This is me learning
how. This is all tested using the `4.19.0-9` kernel.

## Contents

**`hello.c`**: a simple "hello world" module

**`who-connect-me.c`**: a custom netfilter hook to log remote address from TCP SYN packet

**`add-arp-records.c`**: a custom netfilter hook to add arp records to global arp table

**`check-tcp-syncookies.c`**: a custom netfilter hook to check mss from tcp syncookies

**`custom-netlink.c`**: a custom netlink to communicate with Go

**`kprobe_tcp_conn_request`**: a custom kprobe to learn getting argument from kprobing function

**`run-bpf-prog`**: run bpf prog in kernel module, based on [github.com/Asphaltt/iptables-bpf](https://github.com/Asphaltt/iptables-bpf)

~~**`hello-packet.c`**: logs every time your computer receives a packet.
This one could easily be modified to drop packets 50% of the time.~~

~~**`rootkit.c`**: A simple rootkit. [blog post explaining it more](http://jvns.ca/blog/2013/10/08/day-6-i-wrote-a-rootkit/)~~

## Compiling them

I'm running Linux `4.19.0-9`. (run `uname -r`) to find out what you're
using. This almost certainly won't work with a `2.x` kernel, and I
don't know enough. It is unlikely to do any lasting damage to your
computer, but I can't guarantee anything.

```
sudo apt-get install linux-headers-$(uname -r)
```

but I don't remember for sure. If you try this out, I'd love to hear.

To compile them, just run

```
make
```

## Inserting into your kernel (at your own risk!)

sudo insmod hello.ko
dmesg | tail
sudo rmmod hello.ko

insmod ./run-bpf-prog.ko bpf_path=/sys/fs/bpf/iptbpf
ping -c4 114.114.114.114
rmmod `run_bpf_prog`
dmesg | tail



should display the "hello world" message

