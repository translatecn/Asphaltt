#!/bin/bash

CUR_DIR=$(cd `dirname $0` || exit 1; pwd -P)
echo ${CUR_DIR}

set -xe
#git clone https://gitee.com/ls-2018/ebpf.git /tmp/ebpf
#cd /tmp/ebpf/cmd/bpf2go && go install . && cd -

#go install github.com/cilium/ebpf/cmd/bpf2go@latest
mkdir ../bpf || echo 'mkdir success'

cd /data/examples
BPFFILE=../include/libxdp_generated.h

build_bpf() {
    num="$1"
    cat >${BPFFILE} <<EOF
#ifndef __LIBXDP_GENERATED_H_
#define __LIBXDP_GENERATED_H_

#define BITMAP_ARRAY_SIZE ${num}

#endif // __LIBXDP_GENERATED_H_
EOF
    echo `pwd`
    bpf2go -output-dir=../bpf -go-package=bpf "XDPACL${num}" xdp_acl.c --  -D__TARGET_ARCH_x86 -I../include -nostdinc  -Wall -O3
}

main() {
    for x in {8,16,32,64,128,160,256}; do
        build_bpf $x
    done
}


bpf2go -output-dir=../bpf -go-package=bpf                                tcp tcp.c -- -I $BPF_HEADERS


main

bpf2go -output-dir=../bpf -go-package=bpf                                sys sys.bpf.c -- -I $BPF_HEADERS
bpf2go -output-dir=../bpf -go-package=bpf                                xdp xdp.c -- -I $BPF_HEADERS
bpf2go -output-dir=../bpf -go-package=bpf                                tcx tcx.c -- -I $BPF_HEADERS
bpf2go -output-dir=../bpf -go-package=bpf                                tracepoint tracepoint.c -- -I $BPF_HEADERS
bpf2go -output-dir=../bpf -go-package=bpf -type=event                    ringbuffer ringbuffer.c -- -I $BPF_HEADERS
bpf2go -output-dir=../bpf -go-package=bpf                                kprobe_pin kprobe_pin.c -- -I $BPF_HEADERS
#bpf2go -output-dir=../bpf -go-package=bpf -type=event                    tcprtt tcprtt.c -- -I $BPF_HEADERS
bpf2go -output-dir=../bpf -go-package=bpf -type=event                    fentry fentry.c -- -I $BPF_HEADERS
bpf2go -output-dir=../bpf -go-package=bpf                                kprobe kprobe.c -- -I $BPF_HEADERS
bpf2go -output-dir=../bpf -go-package=bpf                                kprobe_percpu kprobe_percpu.c -- -I $BPF_HEADERS
bpf2go -output-dir=../bpf -go-package=bpf                                cgroup_skb cgroup_skb.c -- -I $BPF_HEADERS
bpf2go -output-dir=../bpf -go-package=bpf -type=event -target=amd64      uretprobe uretprobe.c -- -I $BPF_HEADERS
bpf2go -output-dir=../bpf -go-package=bpf                                tc_write tc_write.bpf.c -- -I $BPF_HEADERS
bpf2go -output-dir=../bpf -go-package=bpf                                cgroup_skb cgroup_skb.c -- -I $BPF_HEADERS
