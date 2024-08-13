//go:build ignore
#include "vmlinux.h"
#include "bpf_endian.h"
#include "bpf_helpers.h"


typedef struct event {
    __be32 saddr, daddr;
    __be16 sport, dport;
} __attribute__((packed)) event_t;

// 当前服务器向外发起 tcp 连接、获取接收 tcp 连接时，就将其中的地址端口打印出来。
struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
} events SEC(".maps");

// bpf2go 并没有为全局变量生成对应的 Go struct。

  __be32 filter_daddr;
  __be16 filter_dport;

// const int siz = sizeof(*val);

// __always_inline
static __noinline void handle_new_connection(void *ctx, struct sock *sk) {
    struct sock_common __sk_common;
    event_t ev = {};
    ev.saddr = BPF_CORE_READ(sk, __sk_common.skc_rcv_saddr);
    ev.daddr = BPF_CORE_READ(sk, __sk_common.skc_daddr);
    ev.sport = BPF_CORE_READ(sk, __sk_common.skc_num);
    ev.dport = bpf_ntohs(BPF_CORE_READ(sk, __sk_common.skc_dport));

    if (ev.daddr == filter_daddr && ev.dport == filter_dport) // 使用
        bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &ev, sizeof(ev));
}

SEC("kprobe/tcp_connect")
int k_tcp_connect(struct pt_regs *ctx) {
    struct sock *sk;

    sk = (typeof(sk))PT_REGS_PARM1(ctx);

    handle_new_connection(ctx, sk);

    return 0;
}

SEC("kprobe/inet_csk_complete_hashdance")
int k_icsk_complete_hashdance(struct pt_regs *ctx) {
    struct sock *sk;
    sk = (typeof(sk))PT_REGS_PARM2(ctx);

    handle_new_connection(ctx, sk);

    return 0;
}