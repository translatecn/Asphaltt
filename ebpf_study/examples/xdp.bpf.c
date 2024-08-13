//go:build ignore

#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>

#include "bpf_endian.h"
#include "bpf_helpers.h"

struct ip_data {
    __u32 sip;
    __u32 dip;
    __be16 sport;
    __be16 dport;
};

struct bpf_map_def SEC("maps") allow_ip_maps = {
    .type = BPF_MAP_TYPE_HASH,
    .key_size = sizeof(__u32),
    .value_size = sizeof(__u8),
    .max_entries = 1024,
};

char LICENSE[] SEC("license") = "Dual BSD/GPL";

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 20);
} ip_map SEC(".maps");

SEC("xdp")
int my_pass(struct xdp_md *ctx) {
    void *_data = (void *)(long)ctx->data;
    void *data_end = (void *)(long)ctx->data_end;
    // int pkt_sz = data_end - _data;
    //
    struct ethhdr *eth = _data;
    if ((void *)eth + sizeof(*eth) > data_end) {
        bpf_printk("invalid ethernet header\n");
        return XDP_DROP;
    }

    struct iphdr *ip = _data + sizeof(*ip);
    if ((void *)ip + sizeof(*ip) > data_end) {
        bpf_printk("invalid ip header\n");
        return XDP_DROP;
    }
    struct tcphdr *tcp = _data + sizeof(*tcp);
    if ((void *)tcp + sizeof(*tcp) > data_end) {
        bpf_printk("invalid tcp header\n");
        return XDP_DROP;
    }

    if (ip->protocol != 6) { // 不是 TCP 跳过
        return XDP_PASS;
    }

    struct ip_data *data = NULL;

    data = bpf_ringbuf_reserve(&ip_map, sizeof(struct ip_data), 0); // 环形缓冲区获取一块内存

    if (data) {
        data->sip = bpf_ntohl(ip->saddr); // 网络字节序 转 主机字节序  32位  大小端
        data->dip = bpf_ntohl(ip->daddr);
        data->sport = bpf_ntohs(tcp->source); // 16位
        data->dport = bpf_ntohs(tcp->dest);
        bpf_ringbuf_submit(data, 0);
    }
    __u32 sip = bpf_ntohl(ip->saddr);
    __u8 *allow = bpf_map_lookup_elem(&allow_ip_maps, &sip);

    if (allow && *allow == 1) {
        return XDP_PASS;
    }
    return XDP_DROP;
}