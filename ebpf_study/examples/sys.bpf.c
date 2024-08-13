//go:build ignore

#include "vmlinux.h"
#include "bpf_helpers.h"
#include <linux/limits.h>

struct proc_t {
    __u32 ppid;
    __u32 pid;
    char pname[NAME_MAX];
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 20);
} proc_map SEC(".maps");

static volatile const __u32 XDPACL_DEBUG = 0;

#define bpf_debug_printk(fmt, ...)          \
    do {                                    \
        if (XDPACL_DEBUG)                   \
            bpf_printk(fmt, ##__VA_ARGS__); \
    } while (0)

//SEC("tracepoint/syscalls/sys_enter_execve")

SEC("tracepoint/syscalls/sys_exit_execve")
int handle_tp(void *ctx) {
    struct proc_t *p = NULL;
    p = bpf_ringbuf_reserve(&proc_map, sizeof(*p), 0);
    if (!p) {
        return 0;
    }
    p->pid = bpf_get_current_pid_tgid() >> 32;
    bpf_get_current_comm(&p->pname, sizeof(p->pname));

    struct task_struct *task = (struct task_struct *)bpf_get_current_task();

    if (task) {
        struct task_struct *parent = NULL;
        bpf_probe_read_kernel(&parent, sizeof(*parent), &task->real_parent);
        if (parent) {
            bpf_probe_read_kernel(&p->ppid, sizeof(p->ppid), &parent->pid);
        }
    }

    bpf_ringbuf_submit(p, 0);
    return 0;
}

SEC("kprobe/finish_task_switch")
int handle_sw(struct task_struct *pre) {
    __u32 cur_pid = 0;
    __u32 pre_pid = 0;
    struct task_struct *cur = (struct task_struct *)bpf_get_current_task();
    if (cur) {
        bpf_probe_read_kernel(&cur_pid, sizeof(cur_pid), &cur->pid);
    }
    if (pre) {
        bpf_probe_read_kernel(&pre_pid, sizeof(pre_pid), &pre->pid);
    }
    return 0;
}

struct bash_event {
    __u32 pid;
    u8 line[80];
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 20);
} event_map SEC(".maps");

SEC("uretprobe/readline")
int uretprobe_readline(struct pt_regs *ctx) {
    struct bash_event *event = NULL;
    event = bpf_ringbuf_reserve(&event_map, sizeof(*event), 0);
    if (!event) {
        return 0;
    }
    event->pid = bpf_get_current_pid_tgid() >> 32;
    // PT_REGS_RC 获取函数返回值
    // bpf_probe_read(&event->line, sizeof(event->line), (void *) PT_REGS_RC(ctx)); // TODO 有bug

    bpf_ringbuf_submit(event, 0);
    return 0;
}
