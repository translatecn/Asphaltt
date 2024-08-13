//go:build ignore

#include "bpf_endian.h"
#include "libxdp_acl.h"

char _license[] SEC("license") = "GPL";


SEC("xdp_acl")
int xdp_acl_func_imm(struct xdp_md *ctx) {
    // 该程序就是需要被动态更新的程序。
    return xdp_acl_ipv4(ctx);
}

