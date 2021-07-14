FROM golang:1.16
WORKDIR /bot
COPY . .
ARG test=false
RUN rm -r config/.backup/
RUN bash ./docker-move-configs.sh --test=$TEST
ENTRYPOINT ["go", "test", "./..."]
