```
[root@aps04 ~]# perf list tracepoint |grep sys_enter_write
syscalls:sys_enter_write                           [Tracepoint event]
syscalls:sys_enter_writev                          [Tracepoint event]



root@ubuntu-linux-22-04-desktop:~# cat /sys/kernel/debug/tracing/available_filter_functions | wc -l
60765



4.1 摆脱内核头文件依赖
内核 BTF 信息除了用来做字段重定位之外，还可以用来生成一个大的头文件（"vmlinux.h"），
这个头文件中**包含了所有的内部内核类型，从而避免了依赖系统层面的内核头文件**。
bpftool btf dump file /sys/kernel/btf/vmlinux format c > include/vmlinux.h
<!--  只需要 #include "vmlinux.h"，**也不用再安装 kernel-devel -->

```
- uprobe            挂载在函数进入之前,可以获取到函数的参数值
- uretprobe         挂载在函数返回值之后,可以获取到函数的返回值
- nm /bin/bash   查看一个程序的符号表

xdp 入流量
tc  入出流量
选择网卡 -> chaugnjian duilie  qdisc -> 创建分类class(用于设定宽带级别) -> 创建filter,把流量进行分类,并将包分发到前面定义的class中


# 使用docker0创建一个队列
tc qdisc add dev docker clsact

tc filter add dev docker0 ingress bpf direct-action obj dockertcxdp_bpfel_x86.0 sec.txt

tc filter show dev docker0 ingress 

# 清理掉
tc qdisc del dev docker0 clsact

# 静态hook 点
find /sys/kernel/debug/tracing/events-type d




arp :   广播
ip -> mac地址    封装在了以太网报文里

arping -I etho0 192.168.0.3

arp 欺骗，metallb



- bpf_map_lookup_elem
- bpf_map_update_elem
- bpf_map_delete_elem

- bpf_probe_read  # 如果使用的内核版本还没支持 BPF_PROG_TYPE_TRACING，就必须显式地使用 bpf_probe_read()来读取字段。
- bpf_core_read   # 新的
- bpf_get_current_pid_tgid

- bpf_get_smp_processor_id
- bpf_get_numa_node_id # cat /boot/config-$(uname -r) | grep CONFIG_USE_PERCPU_NUMA_NODE_ID

- bpf_tail_call_static
- bpf_skc_lookup_tcp



```
struct bpf_map_def SEC("maps") kprobe_map = {
	.type        = BPF_MAP_TYPE_PERCPU_ARRAY, // BPF_MAP_TYPE_ARRAY,BPF_MAP_TYPE_PERF_EVENT_ARRAY
	.key_size    = sizeof(u32),
	.value_size  = sizeof(u64),
	.max_entries = 1,
};




```


- BPF_SOCK_OPS_ACTIVE_ESTABLISHED_CB    如果客户端发起连接请求并完成三次握手后的操作符
- BPF_SOCK_OPS_TCP_LISTEN_CB            套接字进入监听状态时的操作符
- BPF_SOCK_OPS_DATA_ACK_CB              数据被确认
- BPF_SOCK_OPS_STATE_CB                 TCP 状态改变


- BPF_SOCK_OPS_RTT_CB_FLAG
- BPF_SOCK_OPS_STATE_CB_FLAG




# 4.3 处理内核版本和配置差异

- https://mp.weixin.qq.com/s?__biz=MzU1MzY4NzQ1OA==&mid=2247493547&idx=1&sn=ab88985daf42faff62f91f2bdb120672&chksm=fbeda766cc9a2e70e503173116ecb1994f26589a82dda5f1b1dbda138725e54470281daef04d&scene=58&subscene=0#rd
- ```
    extern u32 LINUX_KERNEL_VERSION __kconfig;
    extern u32 CONFIG_HZ __kconfig;

    u64 utime_ns;

    if (LINUX_KERNEL_VERSION >= KERNEL_VERSION(4, 11, 0))
        utime_ns = BPF_CORE_READ(task, utime);
    else
        /* convert jiffies to nanoseconds */
        utime_ns = BPF_CORE_READ(task, utime) * (1000000000UL / CONFIG_HZ);

```
- ```
/* up-to-date thread_struct definition matching newer kernels */
struct thread_struct {
    ...
    u64 fsbase;
    ...
};

/* legacy thread_struct definition for <= 4.6 kernels */
struct thread_struct___v46 {   /* ___v46 is a "flavor" part */
    ...
    u64 fs;
    ...
};

extern
int LINUX_KERNEL_VERSION __kconfig;
...

struct thread_struct *thr = ...;
u64 fsbase;
if (LINUX_KERNEL_VERSION > KERNEL_VERSION(4, 6, 0))
    fsbase = BPF_CORE_READ((struct thread_struct___v46 *)thr, fs);
else
    fsbase = BPF_CORE_READ(thr, fsbase);

```


- BPF_PROG_TYPE_SK_REUSEPORT
- BPF_MAP_TYPE_REUSEPORT_SOCKARRAY





-- kprobe 可以跟踪的值 cat /sys/kernel/debug/tracing/available_filter_functions|grep dev_xdp_attach
-- 自定义内核函数 /sys/kernel/btf 一级文件、非vmlinux