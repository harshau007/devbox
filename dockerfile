FROM debian:bullseye-slim

LABEL maintainer="https://github.com/harshau007"
LABEL createdBy="DevBox"

RUN apt-get update && apt-get install -y curl ca-certificates software-properties-common

ARG ADDITIONAL_PACKAGES

ARG ADDITIONAL_PORT

ARG TEMPLATE_NAME

RUN if [ -n "$ADDITIONAL_PACKAGES" ]; then \
    case "$ADDITIONAL_PACKAGES" in \
    "nodelts") \
        curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - && \
        apt-get install -y nodejs ;; \
    "node21") \
        curl -fsSL https://deb.nodesource.com/setup_21.x | bash - && \
        apt-get install -y nodejs ;; \
    "node18") \
        curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
        apt-get install -y nodejs ;; \
    "python") \
        apt-get install -y python3 python3-pip ;; \
    "rust") \
        curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y ;; \
    "go") \
        add-apt-repository ppa:longsleep/golang-backports && \
        apt-get update && \
        apt-get install -y golang-go ;; \
    "java8") \
        apt-get install -y openjdk-8-jdk ;; \
    "java11") \
        apt-get install -y openjdk-11-jdk ;; \
    "java17") \
        apt-get install -y openjdk-17-jdk ;; \
    "java20") \
        apt-get install -y openjdk-20-jdk ;; \
    "java21") \
        apt-get install -y openjdk-21-jdk ;; \
    *) \
        echo $ADDITIONAL_PACKAGES ;; \
    esac; \
fi


RUN curl -fsSL https://code-server.dev/install.sh | sh \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY settings.json /root/.local/share/code-server/User/settings.json

EXPOSE 8080

WORKDIR /home/coder

COPY setup.sh /usr/local/bin/setup.sh

RUN chmod +x /usr/local/bin/setup.sh

# Start code-server
ENTRYPOINT [ "/usr/local/bin/setup.sh" ]