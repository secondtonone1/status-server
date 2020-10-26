FROM alpine:latest
RUN mkdir -p /data/status-server/bin
RUN mkdir -p /data/status-server/log
#ADD ./bin/IMDA /data/bin
ADD  bin/status-server  /data/status-server/bin
COPY conf /data/status-server/conf
COPY src/resources/  /data/status-server/src/resources
EXPOSE 9090
ENV GRPC_ADDR  ""
# Redis数据库的地址: single_redis_host or cluster_redis_host，格式为：redisIMDA:7379
ENV REDIS_HOST ""
RUN echo http://mirrors.aliyun.com/alpine/v3.10/main/ > /etc/apk/repositories && \
    echo http://mirrors.aliyun.com/alpine/v3.10/community/ >> /etc/apk/repositories
RUN apk update && apk upgrade
RUN apk add --no-cache tzdata \
    && ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
ENV TZ Asia/Shanghai
CMD ["/data/status-server/bin/status-server", "--conf","/data/status-server/conf/status_docker.toml"]

