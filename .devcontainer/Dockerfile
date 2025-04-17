ARG VARIANT=1.24-bookworm
FROM golang:${VARIANT}

SHELL ["/bin/bash", "-c"]

# Install X11 development libraries and essential tools
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        libxxf86vm-dev \
        libxinerama-dev \
        libx11-dev \
        libxrandr-dev \
        libxcursor-dev \
        libxi-dev \
        libgl1-mesa-dev \
        xorg-dev \
        libxtst-dev \
        libpng++-dev \
        libasound2-dev \
        xvfb \
        bash \
        curl \
        git \
        vim \
        make \
        build-essential \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Create the user and set up Go environment with proper permissions
RUN groupadd --gid 1000 vscode \
    && useradd --uid 1000 --gid 1000 -m vscode \
    && mkdir -p /home/vscode/go/bin \
    && mkdir -p /go/pkg/mod \
    && mkdir -p /go/pkg/mod/cache/download \
    && mkdir -p /go/pkg/sumdb \
    && chown -R vscode:vscode /home/vscode \
    && chown -R vscode:vscode /go \
    && chmod -R 777 /go/pkg/mod/cache

# Set GOPATH and ensure Go uses the correct permissions
ENV GOPATH=/go
ENV GO111MODULE=on
ENV GOFLAGS=-modcacherw
ENV PATH="/home/vscode/go/bin:${PATH}"

# Switch to vscode user for Go tool installation
USER vscode

# Install Go tools
RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install github.com/golangci/golint/cmd/golint@latest \
    && go install github.com/hajimehoshi/wasmserve@latest

# Set up Xvfb for testing
USER root
RUN echo '#!/bin/bash\n\
if [ "$1" = "test" ]; then\n\
    Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &\n\
    export DISPLAY=:99\n\
    shift\n\
    exec "$@"\n\
else\n\
    exec "$@"\n\
fi' > /usr/local/bin/entrypoint.sh \
    && chmod +x /usr/local/bin/entrypoint.sh

# Set display for X11 forwarding
ENV DISPLAY=host.docker.internal:0.0

# Switch back to vscode user
USER vscode
WORKDIR /workspaces/gimbal

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["sleep", "infinity"]
