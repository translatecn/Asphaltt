FROM registry.cn-hangzhou.aliyuncs.com/acejilam/mygo:v1.22.2
RUN #git clone https://gitee.com/ls-2018/ebpf.git /tmp/ebpf
RUN git clone https://github.com/cilium/ebpf.git /tmp/ebpf
RUN cd /tmp/ebpf/cmd/bpf2go && go install .


# docker build -t registry.cn-hangzhou.aliyuncs.com/acejilam/bpf:v1.22.2 .
# docker push registry.cn-hangzhou.aliyuncs.com/acejilam/bpf:v1.22.2