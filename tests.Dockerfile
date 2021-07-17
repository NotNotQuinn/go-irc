FROM golang:1.16
WORKDIR /bot
COPY . .
RUN cp tests-docker-entrypoint.sh /usr/local/bin

# Install netcat (tbh I dont know if all these commands are needed)
RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y netcat \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN make populator
ENTRYPOINT ["tests-docker-entrypoint.sh"]
