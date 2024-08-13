//go:build ignore

#define MAGIC 0xFEDCBA98

#define ctx_ptr(ctx, mem) (void *)(unsigned long)ctx->mem

struct {
    __uint(type, BPF_MAP_TYPE_PROG_ARRAY);
    __type(key, __u32);
    __type(value, __u32);
    __uint(max_entries, 1);
} xdp_progs SEC(".maps");

SEC("xdp")
int xdp_fn(struct xdp_md *ctx)
{
    __u32 *val;

    // Note: do not bpf_xdp_adjust_meta again.

    void *data_meta = ctx_ptr(ctx, data_meta);
    void *data = ctx_ptr(ctx, data);

    val = (typeof(val))data_meta;
    if ((void *)(val + 1) > data)
        return XDP_PASS;

    if (*val == MAGIC)
        bpf_printk("xdp tailcall\n");

    return XDP_PASS;
}

SEC("xdp")
int xdp_tailcall(struct xdp_md *ctx)
{
    __u32 *val;
    const int siz = sizeof(*val);

    if (bpf_xdp_adjust_meta(ctx, -siz) != 0)
        return XDP_PASS;

    void *data_meta = ctx_ptr(ctx, data_meta);
    void *data = ctx_ptr(ctx, data);

    val = (typeof(val))data_meta;
    if ((void *)(val + 1) > data)
        return XDP_PASS;

    *val = MAGIC;
    bpf_printk("xdp metadata\n");

    bpf_tail_call_static(ctx, &xdp_progs, 0);

    return XDP_PASS;
}

SEC("tc")
int tc_metadata(struct __sk_buff *skb)
{
    void *data = ctx_ptr(skb, data);
    void *data_meta = ctx_ptr(skb, data_meta);

    __u32 *val;
    val = (typeof(val))data_meta;

    if ((void *)(val +1) > data)
        return TC_ACT_OK;

    if (*val == MAGIC)
        bpf_printk("tc metadata\n");

    return TC_ACT_OK;
}

// 利用metadata 在两个 XDP 程序之间传递一个最简单的信息
// xdp_md 中的 data, data_end  指向了   __sk_buff 的data,data_end