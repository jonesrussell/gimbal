FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

# Install X11 development libraries
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
&& apt-get -y install --no-install-recommends \
    libxxf86vm-dev \
    libxinerama-dev \
    libx11-dev \
    libxrandr-dev \
    libxcursor-dev \
    libxi-dev \
    libgl1-mesa-dev \
&& apt-get clean -y && rm -rf /var/lib/apt/lists/*

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
# USER root

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1