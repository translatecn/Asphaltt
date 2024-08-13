/**
 * Copyright 2023 Leon Hwang.
 * SPDX-License-Identifier: Apache-2.0
 */

#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <time.h>
#include <sys/resource.h>
#include <bpf/libbpf.h>
#include "tracepoint.skel.h"

int libbpf_print_fn(enum libbpf_print_level level, const char *format, va_list args)
{
	/* Ignore debug-level libbpf logs */
	if (level > LIBBPF_INFO)
		return 0;
	return vfprintf(stderr, format, args);
}

void bump_memlock_rlimit(void)
{
	struct rlimit rlim_new = {
		.rlim_cur = RLIM_INFINITY,
		.rlim_max = RLIM_INFINITY,
	};

	if (setrlimit(RLIMIT_MEMLOCK, &rlim_new)) {
		fprintf(stderr, "Failed to increase RLIMIT_MEMLOCK limit!\n");
		exit(1);
	}
}

static volatile bool exiting = false;

static void sig_handler(int sig)
{
	exiting = true;
}

struct xdp_errmsg {
	char msg[256];
};

void recv_xdp_errmsg(void *ctx, int cpu, void *data, __u32 data_sz)
{
	struct xdp_errmsg *errmsg = data, *the_errmsg = ctx;
	memcpy(the_errmsg, errmsg, sizeof(*errmsg));
}

int main(int argc, char **argv)
{
	struct perf_buffer *pb = NULL;
	struct perf_buffer_opts pb_opts = {};
	struct xdp_errmsg the_errmsg = {};
	struct tracepoint_bpf *skel;
	int err;

	/* Set up libbpf logging callback */
	libbpf_set_print(libbpf_print_fn);

	/* Bump RLIMIT_MEMLOCK to create BPF maps */
	bump_memlock_rlimit();

	/* Clean handling of Ctrl-C */
	signal(SIGINT, sig_handler);
	signal(SIGTERM, sig_handler);

	skel = tracepoint_bpf__open_and_load();
	if (!skel) {
		fprintf(stderr, "Failed to open and load BPF skeleton\n");
		return 1;
	}

	err = tracepoint_bpf__attach(skel);
	if (err) {
		fprintf(stderr, "Failed to attach BPF skeleton\n");
		goto cleanup;
	}

	pb_opts.sz = sizeof(struct xdp_errmsg);
	pb_opts.sample_period = 1;
	pb = perf_buffer__new(bpf_map__fd(skel->maps.xdp_errmsg_pb), 1, recv_xdp_errmsg, NULL,
			      &the_errmsg, &pb_opts);
	if (libbpf_get_error(pb)) {
		err = -1;
		fprintf(stderr, "Failed to setup perf_buffer\n");
		goto cleanup;
	}

	while (!exiting) {
		err = perf_buffer__poll(pb, 100 /* timeout, ms */);
		/* Ctrl-C will cause -EINTR */
		if (err == -EINTR) {
			err = 0;
			break;
		}
		if (err < 0) {
			fprintf(stderr, "Error polling perf buffer: %d\n", err);
			break;
		}
		if (the_errmsg.msg[0]) {
			printf("Error: %s\n", the_errmsg.msg);
			the_errmsg.msg[0] = '\0';
		}
	}

cleanup:
	perf_buffer__free(pb);
	tracepoint_bpf__destroy(skel);
	return -err;
}