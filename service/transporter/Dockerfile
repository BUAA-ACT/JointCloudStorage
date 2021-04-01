FROM alpine:3.13.4

VOLUME /data
COPY build/bin/transporter /home/jcspan/
CMD chmod +x /home/jcspan/transporter
WORKDIR /home/jcspan
RUN ["./transporter", "-c", "/home/jcspan/config/transporter_config.json"]

EXPOSE 8083