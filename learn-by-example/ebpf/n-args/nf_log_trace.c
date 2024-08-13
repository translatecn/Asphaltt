#define NR_REG_ARGUMENTS 6
#define NR_ARM64_MAX_REG_ARGUMENTS 31

static __inline unsigned long regs_get_kernel_stack_nth_addr(struct pt_regs *regs, unsigned int n) {
    unsigned long *addr = (unsigned long *)regs->sp, retval = 0;

    addr += n;
    return 0 != bpf_probe_read_kernel(&retval, sizeof(retval), addr) ? 0 : retval;
}

static __inline unsigned long regs_get_nth_argument(struct pt_regs *regs, unsigned int n) {
    switch (n) {
        case 0:
            return PT_REGS_PARM1_CORE(regs);
        case 1:
            return PT_REGS_PARM2_CORE(regs);
        case 2:
            return PT_REGS_PARM3_CORE(regs);
        case 3:
            return PT_REGS_PARM4_CORE(regs);
        case 4:
            return PT_REGS_PARM5_CORE(regs);
        case 5:
            return PT_REGS_PARM6_CORE(regs);
        default:
#ifdef __TARGET_ARCH_arm64
            if (n < NR_ARM64_MAX_REG_ARGUMENTS)
                return regs->regs[n];
            else
                return 0;
#elifdef __TARGET_ARCH_x86
            n -= NR_REG_ARGUMENTS - 1;
            return regs_get_kernel_stack_nth_addr(regs, n);
#else
            return 0;
#endif
    }
}

// 用法如下：

SEC("kprobe/nf_log_trace")
int BPF_KPROBE(k_nf_log_trace, struct net *net, u_int8_t pf, unsigned int hooknum, struct sk_buff *skb, struct net_device *in) {
    struct net_device *out;
    char *tablename;
    char *chainname;
    unsigned int rulenum;

    out = (typeof(out))(void *)regs_get_nth_argument(ctx, 5);
    tablename = (typeof(tablename))(void *)regs_get_nth_argument(ctx, 8);
    chainname = (typeof(chainname))(void *)regs_get_nth_argument(ctx, 9);
    rulenum = (typeof(rulenum))regs_get_nth_argument(ctx, 11);

    return __ipt_do_table_trace(ctx, pf, hooknum, skb, in, out, tablename, chainname, rulenum);
}