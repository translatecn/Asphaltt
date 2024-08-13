#include <linux/types.h>

#include "bpf_helpers.h"

#define DEVICE_NUM 256

// m_fdb 存储网卡转发的下一跳关系
BPF_MAP_DEF(m_fdb) = {
    .map_type = BPF_MAP_TYPE_HASH,
    .key_size = sizeof(int),
    .value_size = sizeof(int),
    .max_entries = DEVICE_NUM,
};
BPF_MAP_ADD(m_fdb);

// m_ports 存储网卡信息
BPF_MAP_DEF(m_ports) = {
    .map_type = BPF_MAP_TYPE_DEVMAP,
    .key_size = sizeof(int),
    .value_size = sizeof(int),
    .max_entries = DEVICE_NUM,
};
BPF_MAP_ADD(m_ports);

SEC("xdp")
int xdp_redirect(struct xdp_md *ctx) {
  int ifindex_ingress = ctx->ingress_ifindex;
  int *ifindex_egress = 0;

  ifindex_egress = bpf_map_lookup_elem(&m_fdb, &ifindex_ingress);
  if (ifindex_egress && 0 != *ifindex_egress) {
    return bpf_redirect_map(&m_ports, *ifindex_egress, 0);
  } else {
    return XDP_DROP;
  }
}

SEC("xdp")
int xdp_pass(struct xdp_md *ctx) { return XDP_PASS; }

char _license[] SEC("license") = "GPL";
