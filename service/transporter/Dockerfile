FROM alpine:3.13.4

VOLUME /data
COPY build/bin/transporter /home/jcspan/
RUN chmod +x /home/jcspan/transporter
WORKDIR /home/jcspan
CMD ["./transporter", "-c", "/home/jcspan/config/transporter_config.json"]

EXPOSE 8083