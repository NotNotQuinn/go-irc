FROM golang:1.16

# Move files
WORKDIR /bot
COPY . .
RUN cp tests-docker-entrypoint.sh /usr/local/bin

# To run the bot, set this to anything other than "true"
ARG WB_TEST=true

# lol, for some stats. i guess.
ENV WB_DOCKER=true

# Install netcat (tbh I dont know if all these commands are needed)
RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y netcat \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create binaries
RUN make all
ENTRYPOINT ["tests-docker-entrypoint.sh"]
