FROM scratch
MAINTAINER cblomart@gmail.com

COPY ../../releases/linux/amd64/zwaymqtt /zwaymqtt
RUN chmod 777 /zwaymqtt

ENTRYPOINT [ "/zwaymqtt" ] 