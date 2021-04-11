FROM harbor.sumblog.cn/pub/alpine:3.13.4

VOLUME /data
COPY build/bin/scheduler /home/jcspan/
RUN chmod +x /home/jcspan/scheduler
WORKDIR /home/jcspan
CMD ["cd", "/home/jcspan"]
CMD ["ls", "-a"]
CMD ["./scheduler", "-addr=:8082", "-cid=aliyun-beijing", "-env=dev", "-heartbeat=10s", "-reschedule=60s"]

EXPOSE 8082
