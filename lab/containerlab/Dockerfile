FROM ubuntu:24.04

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        iproute2 \
        openssh-server \
    && rm -rf /var/lib/apt/lists/*

CMD ["tail", "-f", "/dev/null"]
