# 检查是否运行支持 iptables-bpf


```
iptables -I OUTPUT -m bpf --object-pinned /sys/fs/bpf/iptbpf -j DROP

### 不支持需要重新编译
git clone git://git.netfilter.org/iptables.git
cd iptables
bash autogen.sh
apt install -y libpcap-dev
./configure --enable-bpf-compiler --disable-nftables # disable nftables 是为了快速安装一个能用 bpf 的 iptables
make -j4
make install
```

