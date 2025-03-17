ARG VARIANT=1.24-bookworm
FROM golang:${VARIANT}

SHELL ["/bin/bash", "-c"]

# Install X11 development libraries, bash, and Node.js prerequisites
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
&& apt-get -y install --no-install-recommends \
    libxxf86vm-dev \
    libxinerama-dev \
    libx11-dev \
    libxrandr-dev \
    libxcursor-dev \
    libxi-dev \
    libgl1-mesa-dev \
    bash \
    curl \
    gnupg \
    zsh \
&& apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Install Node.js LTS
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get update \
    && apt-get install -y nodejs \
    && npm install -g npm@latest

# Create the user and set up Go environment
RUN groupadd --gid 1000 vscode \
    && useradd --uid 1000 --gid 1000 -m vscode \
    && mkdir -p /home/vscode/go/bin \
    && mkdir -p /go/pkg/mod \
    && mkdir -p /go/bin \
    && mkdir -p /go/pkg/mod/cache \
    && mkdir -p /go/pkg/mod/cache/download \
    && mkdir -p /go/pkg/sumdb \
    && chmod -R 777 /go \
    && find /go -type d -exec chown vscode:vscode {} \; \
    && find /go -type f -exec chown vscode:vscode {} \; \
    && chown -R vscode:vscode /home/vscode

# Set GOPATH and ensure Go uses the correct permissions
ENV GOPATH=/go
ENV GO111MODULE=on
ENV GOFLAGS=-modcacherw

# Switch to vscode user to install golang tools
USER vscode
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    go install github.com/vektra/mockery/v2@v2.46.3 && \
    go install github.com/air-verse/air@latest && \
    go install github.com/spf13/cobra-cli@latest && \
    go install github.com/hajimehoshi/wasmserve@latest

# Ensure the migrate binary is in the PATH for vscode user
ENV PATH="/home/vscode/go/bin:${PATH}"

# Switch back to root for any further operations
USER root

# Create zsh completion directory and fix permissions
RUN mkdir -p /usr/local/share/zsh/site-functions && \
    touch /usr/local/share/zsh/site-functions/_task && \
    chmod -R 755 /usr/local/share/zsh && \
    chmod 644 /usr/local/share/zsh/site-functions/_task

# Switch back to vscode user as final user
USER vscode

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "npm install -g <your-package-here>"