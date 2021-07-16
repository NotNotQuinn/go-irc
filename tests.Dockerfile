FROM golang:1.16
WORKDIR /bot
COPY . .
RUN cp tests-docker-entrypoint.sh /usr/local/bin
ENTRYPOINT ["tests-docker-entrypoint.sh"]
