FROM ubuntu:24.04

RUN apt-get update \
    -o Acquire::Retries=5 \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        -o Acquire::Retries=5 \
        iproute2 \
        openssh-server \
    && rm -rf /var/lib/apt/lists/*

CMD ["tail", "-f", "/dev/null"]
