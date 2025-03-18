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
    xorg-dev \
    libxtst-dev \
    libpng++-dev \
    bash \
    curl \
    gnupg \
    zsh \
    git \
    wget \
    vim \
    make \
    build-essential \
    libxext-dev \
    libxrender-dev \
    x11-apps \
    xvfb \
    libasound2-dev \
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

# Ensure the migrate binary is in the PATH for vscode user
ENV PATH="/home/vscode/go/bin:${PATH}"

# Switch back to root for any further operations
USER root

# Create zsh completion directory and fix permissions
RUN mkdir -p /usr/local/share/zsh/site-functions && \
    touch /usr/local/share/zsh/site-functions/_task && \
    chown -R root:root /usr/local/share/zsh && \
    chmod -R 755 /usr/local/share/zsh && \
    chmod 644 /usr/local/share/zsh/site-functions/_task && \
    # Ensure Go directories have correct permissions after all operations
    find /go -type d -exec chown vscode:vscode {} \; && \
    find /go -type f -exec chown vscode:vscode {} \;

# Switch back to vscode user as final user
USER vscode

# Install Go tools
RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
    && go install github.com/cosmtrek/air@latest \
    && go install github.com/hajimehoshi/wasmserve@latest

# Set environment variables
ENV DISPLAY=host.docker.internal:0.0
ENV XAUTHORITY=/home/vscode/.Xauthority

# Create non-root user
RUN echo "vscode ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers.d/vscode \
    && chmod 0440 /etc/sudoers.d/vscode

WORKDIR /workspaces/gimbal

# Create a script to start Xvfb and run tests
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

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["sleep", "infinity"]
