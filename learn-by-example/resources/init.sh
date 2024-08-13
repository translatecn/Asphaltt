set -x

touch ~/.hushlogin

ARCH=$(arch | sed s/aarch64/arm64/ | sed s/x86_64/amd64/)

function install_python() {
	apt install python3-pip -y
	pip3 config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
	pip3 install bcc pytest
	ln -s /usr/bin/python3 /usr/bin/python
}

function init_repo() {
	if [ "$ARCH" == "amd64" ]; then
		cat >/etc/apt/sources.list <<EOF
deb http://mirrors.huaweicloud.com/repository/ubuntu/ jammy main restricted
deb http://mirrors.huaweicloud.com/repository/ubuntu/ jammy-updates main restricted
deb http://mirrors.huaweicloud.com/repository/ubuntu/ jammy universe
deb http://mirrors.huaweicloud.com/repository/ubuntu/ jammy-updates universe
deb http://mirrors.huaweicloud.com/repository/ubuntu/ jammy multiverse
deb http://mirrors.huaweicloud.com/repository/ubuntu/ jammy-updates multiverse
deb http://mirrors.huaweicloud.com/repository/ubuntu/ jammy-backports main restricted universe multiverse
deb http://mirrors.huaweicloud.com/repository/ubuntu/ jammy-security main restricted
deb http://mirrors.huaweicloud.com/repository/ubuntu/ jammy-security universe
deb http://mirrors.huaweicloud.com/repository/ubuntu/ jammy-security multiverse
EOF
	fi

	if [ "$ARCH" == "arm64" ]; then
		cat >/etc/apt/sources.list <<EOF
deb http://repo.huaweicloud.com/ubuntu-ports/ jammy main restricted universe multiverse
deb http://repo.huaweicloud.com/ubuntu-ports/ jammy-security main restricted universe multiverse
deb http://repo.huaweicloud.com/ubuntu-ports/ jammy-updates main restricted universe multiverse
deb http://repo.huaweicloud.com/ubuntu-ports/ jammy-backports main restricted universe multiverse
EOF
	fi
	echo 'nameserver 114.114.114.114' >/etc/resolv.conf
	echo -e 'root\nroot\n' | passwd root
	apt clean all
	apt clean all
	apt clean all
	apt-get update -y
}

function install_bcc() {

	cat <<EOF >>/etc/profile
export PATH=\$PATH:/usr/share/bcc/tools:/usr/lib/llvm-14/bin
EOF
	source /etc/profile

	# 二进制包
	apt-get install -y make clang llvm libelf-dev libbpf-dev bpfcc-tools libbpfcc-dev bpfcc-tools linux-headers-$(uname -r)
	# 有一个安装成功就行
	apt-get install -y linux-tools-$(uname -r)
	apt-get install -y linux-tools-generic

	# 源码编译
	#  apt install -y zip bison build-essential cmake flex git libedit-dev \
	#    libllvm14 llvm-14-dev libclang-14-dev python3 zlib1g-dev libelf-dev libfl-dev python3-setuptools \
	#    liblzma-dev libdebuginfod-dev arping netperf iperf
	# export https_proxy=http://192.168.31.50:7890 http_proxy=http://192.168.31.50:7890 all_proxy=socks5://192.168.31.50:7890
	# export https_proxy=http://10.230.205.190:7890 http_proxy=http://10.230.205.190:7890 all_proxy=socks5://10.230.205.190:7890
	# git clone https://github.com/iovisor/bcc.git
	# mkdir bcc/build; cd bcc/build
	# cmake ..
	# unset https_proxy && unset http_proxy && unset all_proxy
	# make -j 16
	# sudo make install
	# cmake -DPYTHON_CMD=python3 .. # build python3 binding
	# pushd src/python/
	# make -j 16
	# sudo make install
	# popd

	# # 测试
	# /usr/share/bcc/tools/biolatency

	# libbpf-go的依赖
	apt install libbpf-dev -y

	#	export https_proxy=http://192.168.31.50:7890 http_proxy=http://192.168.31.50:7890 all_proxy=socks5://192.168.31.50:7890
	#	git clone --depth 1 https://github.com/libbpf/libbpf
	#	cd libbpf/src
	#	make install
	#	unset https_proxy && unset http_proxy && unset all_proxy

	if [ "$ARCH" == "amd64" ]; then
		ln -s /usr/include/x86_64-linux-gnu/asm /usr/include/asm
	fi
	if [ "$ARCH" == "arm64" ]; then
		ln -s /usr/include/aarch64-linux-gnu/asm /usr/include/asm
	fi

}

function install_go() {
	rm -rf /usr/local/go*
	rm -rf ./go*
	yum install wget vim gcc -y || apt install wget vim gcc -y

	version=1.22.0
	mkdir /usr/local/go$version

	wget https://golang.google.cn/dl/go$version.linux-$ARCH.tar.gz
	tar -xvf go$version.linux-$ARCH.tar.gz -C /usr/local/go$version --strip-components 1
	rm -rf go$version.linux-$ARCH.tar.gz
	mkdir -p ~/.go/{bin,src,pkg}
	chmod -R 777 /usr/local/go$version
	cat <<EOF >>/etc/profile
export GOROOT="/usr/local/go$version"
export GOPATH=\$HOME/.go  #工作地址路径
export GOBIN=\$GOROOT/bin
export PATH=\$PATH:\$GOBIN
EOF
	source /etc/profile
	go version
	go env
	go env -w GO111MODULE=on
	go env -w GOPROXY=https://goproxy.cn,direct
	go env -w GOFLAGS="-buildvcs=false"
	go env -w CGO_ENABLED="1"

	go install github.com/trzsz/trzsz-go/cmd/...@latest
	go install github.com/go-delve/delve/cmd/dlv@master

}

function install_clang() {
	# https://zhuanlan.zhihu.com/p/592334845

	# 自动安装
	# sudo bash -c "$(wget -O - https://apt.llvm.org/llvm.sh)"
	lsb_release -a
	cat >>/etc/apt/sources.list <<EOF
deb http://apt.llvm.org/focal/ llvm-toolchain-focal-15 main
deb-src http://apt.llvm.org/focal/ llvm-toolchain-focal-15 main
EOF
	sudo apt update
	sudo apt upgrade
	sudo apt install -y clang-15 lldb-15 lld-15 libc++-15-dev libc++abi-15-dev

	dpkg -l | grep clang | awk '{print $2}'
	clang --version
	clang++ --version
}

function prepare_ebpf() {
	wget https://aka.pw/bpf-ecli -O /usr/bin/ecli && chmod +x /usr/bin/ecli
	wget https://github.com/eunomia-bpf/eunomia-bpf/releases/download/v1.0.21/ecc -O /usr/bin/ecc && chmod +x /usr/bin/ecc
}

init_repo
install_python
install_bcc
install_go
