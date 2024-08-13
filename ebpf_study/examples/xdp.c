//go:build ignore

#include "bpf_endian.h"
#include "common.h"

char LICENSE[] SEC("license") = "Dual MIT/GPL";

#define MAX_MAP_ENTRIES 16

/* Define an LRU hash map for storing packet count by source IPv4 address */
struct {
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __uint(max_entries, MAX_MAP_ENTRIES);
    __type(key, __u32);   // source IPv4 address
    __type(value, __u32); // packet count
} xdp_stats_map SEC(".maps");

/*
Attempt to parse the IPv4 source address from the packet.
Returns 0 if there is no IPv4 header field; otherwise returns non-zero.
*/
static __always_inline int parse_ip_src_addr(struct xdp_md *ctx, __u32 *ip_src_addr) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;

    // First, parse the ethernet header.
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end) {
        return 0;
    }

    if (eth->h_proto != bpf_htons(ETH_P_IP)) {
        // The protocol is not IPv4, so we can't parse an IPv4 source address.
        return 0;
    }

    // Then parse the IP header.
    struct iphdr *ip = (void *)(eth + 1);
    if ((void *)(ip + 1) > data_end) {
        return 0;
    }

    // Return the source IP address in network byte order.
    *ip_src_addr = (__u32)(ip->saddr);
    return 1;
}

SEC("xdp")
int xdp_prog_func(struct xdp_md *ctx) {
    __u32 ip;
    if (!parse_ip_src_addr(ctx, &ip)) {
        // Not an IPv4 packet, so don't count it.
        goto done;
    }

    __u32 *pkt_count = bpf_map_lookup_elem(&xdp_stats_map, &ip);
    if (!pkt_count) {
        // No entry in the map for this IP address yet, so set the initial value to 1.
        __u32 init_pkt_count = 1;
        bpf_map_update_elem(&xdp_stats_map, &ip, &init_pkt_count, BPF_ANY);
    }
    else {
        // Entry already exists for this IP address,
        // so increment it atomically using an LLVM built-in.
        __sync_fetch_and_add(pkt_count, 1);
    }

done:
    // Try changing this to XDP_DROP and see what happens!
    return XDP_PASS;
}



SEC("xdp")
int my_pass(struct xdp_md *ctx) {
    void *data = (void *)(long)ctx->data;
    void *data_end = (void *)(long)ctx->data_end;
    int pkt_sz = data_end - data;

    struct ethhdr *eth = data;
    if ((void *)eth + sizeof(*eth) > data_end) {
        bpf_printk("invalid ethernet header\n");
        return XDP_DROP;
    }

    struct iphdr *ip = data + sizeof(*ip);
    if ((void *)ip + sizeof(*ip) > data_end) {
        bpf_printk("invalid ip header\n");
        return XDP_DROP;
    }

    unsigned int src_ip = ip->saddr;
    unsigned char bytes[4];
    bytes[0] = (src_ip >> 0) & 0xFF;
    bytes[1] = (src_ip >> 8) & 0xFF;
    bytes[2] = (src_ip >> 16) & 0xFF;
    bytes[3] = (src_ip >> 24) & 0xFF;

    bpf_printk("packet size is %d, protocol is %d, ip is %d.%d.%d.%d\n", pkt_sz, ip->protocol, bytes[0], bytes[1], bytes[2], bytes[3]);
    return XDP_PASS;
}

// 从 XDP 传递 metadata 到 tc
//SEC("xdp")
//int xdp_tailcall(struct xdp_md *ctx)
//{
//    __u32 *val;
//    const int siz = sizeof(*val);
//
//    if (bpf_xdp_adjust_meta(ctx, -siz) != 0)
//        return XDP_PASS;
//
//    void *data_meta = ctx_ptr(ctx, data_meta);
//    void *data = ctx_ptr(ctx, data);
//
//    val = (typeof(val))data_meta;
//    if ((void *)(val + 1) > data)
//        return XDP_PASS;
//
//    *val = MAGIC;
//    bpf_printk("xdp metadata\n");
//
//    bpf_tail_call_static(ctx, &xdp_progs, 0);
//
//    return XDP_PASS;
//}

//SEC("tc")
//int tc_metadata(struct __sk_buff *skb)
//{
//    void *data = ctx_ptr(skb, data);
//    void *data_meta = ctx_ptr(skb, data_meta);
//
//    __u32 *val;
//    val = (typeof(val))data_meta;
//
//    if ((void *)(val +1) > data)
//        return TC_ACT_OK;
//
//    if (*val == MAGIC)
//        bpf_printk("tc metadata\n");
//
//    return TC_ACT_OK;
//}



// XDP_DROP 丢弃网络包--常用DDos防范
// XDP_REDIRECT 数据包转发到不同的网络接口
// XDP_TX 将数据包转发到接收它的同一网络接口
// XDP_PASS 正常流转,流转之前可以做修改
// XDP_ABORTED 丢弃并抛出异常

// ip link set dev docker0 xdp obj test.bpf.o sec xdp verbose
// ip link set dev docker0 xdp off

// 启动两个容器，分别从宿主机 和另一个容器访问
