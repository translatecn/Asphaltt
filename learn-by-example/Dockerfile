FROM registry.cn-hangzhou.aliyuncs.com/acejilam/ebpm
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/root/.eunomia/bin
#FROM registry.cn-hangzhou.aliyuncs.com/acejilam/mygo:v1.22.2
COPY resources/sources.list /etc/apt/sources.list
COPY resources/resolv.conf  /etc/apt/resolv.conf
RUN apt clean all && apt update -y
RUN apt install curl clang-format cmake -y

COPY resources/init.sh /tmp/init.sh
RUN bash /tmp/init.sh

ENTRYPOINT ["bash"]
