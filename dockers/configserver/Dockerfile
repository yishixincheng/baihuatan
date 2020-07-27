FROM alpine:latest

COPY jdk-14.0.2_linux-x64_bin.tar.gz /opt/jdk-14.0.2_linux-x64_bin.tar.gz

RUN build_pkgs="build-base linux-headers openssl-dev pcre-dev wget zlib-dev" \
  && runtime_pkgs="openssl pcre zlib" \
  && apk --update add ${build_pkgs} ${runtime_pkgs} \
  && mkdir -p /opt/java \
  && tar xzf /opt/jdk-14.0.2_linux-x64_bin.tar.gz -C /opt/java \
  && rm -rf /opt/jdk-14.0.2_linux-x64_bin.tar.gz
  
#设置环境变量
ENV JAVA_HOME /opt/java/jdk-14.0.2
ENV PATH JAVA_HOME/bin:PATH

EXPOSE 8050

VOLUME ["/javaproject"]

# 执行命令
CMD ["nohup java" "-jar", "/javaproject/config-server.jar", "&"]

