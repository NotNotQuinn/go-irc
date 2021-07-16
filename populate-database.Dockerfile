FROM golang:1.16
WORKDIR /bot
COPY . .
SHELL [ "bash", "-c" ]
RUN make populator
ENTRYPOINT [ "./populator-entrypoint.sh" ]
