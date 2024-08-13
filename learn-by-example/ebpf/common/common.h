

// 判断一个网络包是不是 IPv4 的 TCP 包的目的端口是不是 8080
static __always_inline bool __is_ipv4_tcp_udp(struct xdp_md *xdp) {
    void *data_end = (void *)(long)xdp->data_end;
    void *data = (void *)(long)xdp->data;
    struct ethhdr *eth;
    struct iphdr *iph;
    struct tcphdr *tcph;

    eth = data;
    iph = data + sizeof(*eth);
    tcph = data + sizeof(*eth) + sizeof(*iph);

    // if ((void *)eth + sizeof(*eth) > data_end)
    //     return false;
    // if ((void *)iph + sizeof(*iph) > data_end)
    //     return false;
    //只需要检查偏移量最大的那一次即可。
    if ((void *)tcph + sizeof(*tcph) > data_end)
        return false;
    if (eth->h_proto != bpf_htons(ETH_P_IP) || iph->protocol != IPPROTO_TCP)
        return false;

    return tcph->dest == bpf_htons(8080);
}