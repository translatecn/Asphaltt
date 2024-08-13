// 用于暂存到map的struct
struct temp_key_t {
    u32 tgid;
    u32 pid;
};

struct temp_value_t {
    u64 start_time; // 触发当前 PID 切出的时间。
    u64 user_stack_id; // 当前 PID 的用户态栈信息的 ID 值。
    u64 kernel_stack_id;// 当前 PID 的内核态栈信息的 ID 值。
    u8 comm[16]; // PID 的名称。
};

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1024);
    __type(key, struct temp_key_t);
    __type(value, struct temp_value_t);
} temp_pid_status SEC(".maps");


// 尝试记录 OFF CPU 的开始时间
inline void try_record_start(void *ctx, u32 prev_pid, u32 prev_tgid) {
    if (prev_tgid == 0 || prev_pid == 0 || prev_tgid != listen_tgid) {
        return;
    }
    // TODO
}


SEC("tp_btf/sched_switch")
int BPF_PROG(sched_switch, bool preempt, struct task_struct *prev, struct task_struct *next) {
    pid_t prev_pid = prev->pid;
    pid_t prev_tgid = prev->tgid;

    pid_t next_pid = next->pid;
    pid_t next_tgid = next->tgid;

    return 0;
}