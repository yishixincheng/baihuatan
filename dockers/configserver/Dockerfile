FROM alpine:latest

RUN build_pkgs="build-base linux-headers openssl-dev pcre-dev wget zlib-dev" \
  && runtime_pkgs="openssl pcre zlib" \
  && apk --update add ${build_pkgs} ${runtime_pkgs}

COPY jdk-14.0.2_linux-x64_bin.tar.gz /opt/jdk-14.0.2_linux-x64_bin.tar.gz
RUN mkdir -p /opt/java

#解压
RUN tar -xzvf /opt/jdk-14.0.2_linux-x64_bin.tar.gz -C /opt/java
RUN rm -rf /opt/jdk-14.0.2_linux-x64_bin.tar.gz
#设置环境变量
ENV JAVA_HOME /opt/java/jdk-14.0.2
ENV PATH JAVA_HOME/bin:PATH

EXPOSE 8050

VOLUME ["/javaproject"]

# 执行命令
CMD ["nohup java" "-jar", "/javaproject/config-server.jar", "&"]

