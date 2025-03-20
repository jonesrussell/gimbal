ARG VARIANT=1.24-bookworm
FROM golang:${VARIANT}

SHELL ["/bin/bash", "-c"]

# Install necessary tools and libraries
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive && apt-get -y install --no-install-recommends \
    libxxf86vm-dev libxinerama-dev libx11-dev libxrandr-dev libxcursor-dev libxi-dev \
    libgl1-mesa-dev xorg-dev libxtst-dev libpng++-dev bash curl gnupg zsh git wget vim \
    make build-essential libxext-dev libxrender-dev x11-apps xvfb libasound2-dev sudo \
    && curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get update && apt-get install -y nodejs \
    && npm install -g npm@latest \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Create the user and set up the environment
RUN groupadd --gid 1000 vscode && \
    useradd --uid 1000 --gid 1000 -m vscode && \
    mkdir -p /home/vscode/go/bin /go/pkg/mod /go/bin /go/pkg/mod/cache /go/pkg/sumdb && \
    chmod -R 777 /go && \
    chown -R vscode:vscode /home/vscode /go

# Set environment variables
ENV GOPATH=/go \
    GO111MODULE=on \
    GOFLAGS=-modcacherw \
    DISPLAY=host.docker.internal:0.0 \
    XAUTHORITY=/home/vscode/.Xauthority \
    PATH="/home/vscode/go/bin:${PATH}"

# Switch to root for further setup
USER root

# Configure sudo and zsh
RUN mkdir -p /usr/local/share/zsh/site-functions && \
    touch /usr/local/share/zsh/site-functions/_task && \
    chmod -R 755 /usr/local/share/zsh && \
    echo "vscode ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/vscode && \
    chmod 0440 /etc/sudoers.d/vscode

# Install Go tools
RUN go install golang.org/x/tools/gopls@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install github.com/hajimehoshi/wasmserve@latest

# Create a script to start Xvfb and run tests
RUN echo '#!/bin/bash\n\
if [ "$1" = "test" ]; then\n\
    Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &\n\
    export DISPLAY=:99\n\
    shift\n\
    exec "$@"\n\
else\n\
    exec "$@"\n\
fi' | tee /usr/local/bin/entrypoint.sh > /dev/null && \
    chmod +x /usr/local/bin/entrypoint.sh

# Use the 'vscode' user as default
USER vscode

WORKDIR /workspaces/gimbal

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["sleep", "infinity"]
