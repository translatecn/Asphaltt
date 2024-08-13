/**
 * Copyright 2023 Leon Hwang.
 * SPDX-License-Identifier: MIT
 */
//go:build ignore

#include "bpf_all.h"

#include "lib_tp_msg.h"

struct netlink_extack_error_ctx {
    unsigned long unused;

    /*
     * bpf does not support tracepoint __data_loc directly.
     *
     * Actually, this field is a 32 bit integer whose value encodes
     * information on where to find the actual data. The first 2 bytes is
     * the size of the data. The last 2 bytes is the offset from the start
     * of the tracepoint struct where the data begins.
     * -- https://github.com/iovisor/bpftrace/pull/1542
     */
    __u32 msg; // __data_loc char[] msg;
};

//确定 tracepoint 的 ctx 结构体的其它字段信息
//方式1
// cat /sys/kernel/debug/tracing/events/netlink/netlink_extack/format
//方式2
// bpftrace -lv 'tracepoint:netlink:netlink_extack'

SEC("tp/netlink/netlink_extack")
int tp__netlink_extack(struct netlink_extack_error_ctx *ctx) {
    char *msg = (void *)(__u64)((void *)ctx + (__u64)((ctx->msg) & 0xFFFF)); // 获取其中的 msg

    __output_msg(ctx, msg, PROBE_TYPE_DEFAULT, 0);

    return 0;
}