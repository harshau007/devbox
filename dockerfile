FROM ubuntu:20.04

LABEL maintainer="https://github.com/harshau007"
LABEL createdBy="devcraft"

RUN apt-get update && apt-get install -y curl ca-certificates software-properties-common

ARG ADDITIONAL_PACKAGES

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
            *) \
                apt-get install -y $ADDITIONAL_PACKAGES ;; \
        esac; \
    fi

RUN curl -fsSL https://code-server.dev/install.sh | sh

COPY settings.json /root/.local/share/code-server/User/settings.json

EXPOSE 8080

WORKDIR /home/coder

# Start code-server
CMD ["code-server", "--bind-addr", "0.0.0.0:8080", ".", "--auth", "none"]