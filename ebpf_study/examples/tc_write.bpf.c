//go:build ignore
#include <linux/bpf.h>
#include <linux/limits.h>

#include "bpf_helpers.h"

typedef unsigned long long pid_t;

char LICENSE[] SEC("license") = "Dual BSD/GPL";

//struct bpf_map_def SEC("maps") my_bpf_map = {
//    .type = BPF_MAP_TYPE_PERF_EVENT_ARRAY,
//    .key_size = sizeof(int),
//    .value_size = sizeof(int),
//    .max_entries = 1024,
//};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 20);
} log_map SEC(".maps");

// 唤醒缓冲区 >=5.8
// 解决了BOF性能缓冲区的内存效率和事件重排序问题
// 多生产者但消费者队列，可以同时在多个CPU之间安全共享
// 内核-> 用户空间的 优先选择

struct data_t {
    pid_t pid;
    char comm[NAME_MAX];
};

SEC("tracepoint/syscalls/sys_enter_write")
int handle_tp(void *ctx) {
    struct data_t *data = NULL;

    data = bpf_ringbuf_reserve(&log_map, sizeof(struct data_t), 0); // 环形缓冲区获取一块内存

    if (data) {
        data->pid = bpf_get_current_pid_tgid() >> 32;
        bpf_get_current_comm(&data->comm, sizeof(data->comm));
        // bpf_perf_event_output(ctx, &my_bpf_map, 0, &data, sizeof(*data));
        bpf_ringbuf_submit(data, 0);
    }
    return 0;
}
