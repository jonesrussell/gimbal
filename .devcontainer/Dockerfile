ARG VARIANT=1-1-bookworm
FROM mcr.microsoft.com/devcontainers/go:${VARIANT}

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

# Set up Xvfb for testing
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
