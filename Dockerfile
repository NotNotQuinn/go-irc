FROM golang:1.16
WORKDIR /bot
COPY . .
ENTRYPOINT ["go", "test", "./..."]
