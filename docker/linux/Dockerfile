FROM scratch
MAINTAINER cblomart@gmail.com
COPY ./zwaymqtt /zwaymqtt
COPY ./ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/zwaymqtt"]