FROM debian:buster-slim

LABEL maintainer="yishixincheng@163.com"

COPY jdk-14.0.2_linux-x64_bin.tar.gz /opt/jdk-14.0.2_linux-x64_bin.tar.gz

RUN build_pkgs="wget" \
  && runtime_pkgs="ca-certificates" \
  && apt-get update \
  && apt-get install --no-install-recommends --no-install-suggests -y ${build_pkgs} ${runtime_pkgs} \
  && mkdir -p /opt/java \
  && tar -xzvf /opt/jdk-14.0.2_linux-x64_bin.tar.gz -C /opt/java \
  && rm -rf /opt/jdk-14.0.2_linux-x64_bin.tar.gz

#设置环境变量
ENV JAVA_HOME /opt/java/jdk-14.0.2
ENV PATH JAVA_HOME/bin:PATH

EXPOSE 8050

VOLUME ["/javaproject"]

# 执行命令
CMD ["nohup java" "-jar", "/javaproject/config-server.jar", "&"]

