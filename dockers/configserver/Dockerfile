FROM debian:buster-slim

COPY jdk-14.0.2_linux-x64_bin.tar.gz /opt/jdk-14.0.2_linux-x64_bin.tar.gz
RUN mkdir -p /opt/java

#解压
RUN tar -xvf /opt/jdk-14.0.2_linux-x64_bin.tar.gz -C /opt/java
RUN rm -rf /opt/jdk-14.0.2_linux-x64_bin.tar.gz
#设置环境变量
ENV JAVA_HOME=/opt/java/jdk-14.0.2
ENV PATH="$JAVA_HOME/bin:${PATH}"

EXPOSE 8050

VOLUME ["/javaproject"]

# 执行命令
CMD ["java" "-jar", "/javaproject/config-server.jar", "&"]
