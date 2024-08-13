/**
 * Copyright 2023 Leon Hwang.
 * SPDX-License-Identifier: Apache-2.0
 */

#include <linux/bpf.h>
#include <linux/ptrace.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

char LICENSE[] SEC("license") = "GPL";

struct xdp_errmsg {
	char msg[256];
};

struct {
	__uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
	__type(key, int);
	__type(value, struct xdp_errmsg);
} xdp_errmsg_pb SEC(".maps");

struct xdp_attach_error_ctx {
	unsigned long unused;

	const char *msg;
};

SEC("tp/xdp/bpf_xdp_link_attach")
int tracepoint__xdp__bpf_xdp_link_attach(struct xdp_attach_error_ctx *ctx)
{
	struct xdp_errmsg errmsg;
	bpf_probe_read_kernel_str(&errmsg.msg, sizeof(errmsg.msg), ctx->msg);
	bpf_perf_event_output(ctx, &xdp_errmsg_pb, BPF_F_CURRENT_CPU, &errmsg, sizeof(errmsg));
	return 0;
}