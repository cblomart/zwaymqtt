FROM scratch
MAINTAINER cblomart@gmail.com
COPY ./zwaymqtt /zwaymqtt
COPY /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/zwaymqtt"]