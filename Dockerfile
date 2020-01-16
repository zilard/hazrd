FROM ubuntu:16.04

ADD hazrd /bin

EXPOSE 8080

CMD ["hazrd"]
