FROM scratch
MAINTAINER cblomart@gmail.com

COPY releases/linux/arm/zwaymqtt /zwaymqtt
RUN chmod 777 /zwaymqtt

ENTRYPOINT [ "/zwaymqtt" ] 