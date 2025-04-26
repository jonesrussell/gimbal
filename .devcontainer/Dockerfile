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
        libgl1-mesa-glx \
        mesa-utils \
        xorg-dev \
        libxtst-dev \
        libpng++-dev \
        libasound2-dev \
        xvfb \
        xauth \
        x11-xserver-utils \
        bash \
        curl \
        git \
        vim \
        make \
        build-essential \
        libxkbfile-dev \
        libsecret-1-dev \
        libxss1 \
        libglu1-mesa \
        libnss3 \
        libatk1.0-0 \
        libatk-bridge2.0-0 \
        libgdk-pixbuf2.0-0 \
        libgtk-3-0 \
        libgbm-dev \
        x11-apps \
        dbus-x11 \
        libxcb1-dev \
        libx11-xcb-dev \
        libxcb-keysyms1-dev \
        libxcb-util0-dev \
        libxcb-icccm4-dev \
        libxcb-render-util0-dev \
        libxcb-xkb-dev \
        libxkbcommon-dev \
        libxkbcommon-x11-dev \
        vulkan-tools \
    && rm -rf /var/lib/apt/lists/*

# Set up Xvfb for testing
RUN echo '#!/bin/bash\n\
if [ "$1" = "test" ]; then\n\
    # Set up fake X server with authentication\n\
    export XAUTHORITY=/tmp/.Xauthority\n\
    touch $XAUTHORITY\n\
    xauth generate :99 . trusted\n\
    Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &\n\
    export DISPLAY=:99\n\
    shift\n\
    exec "$@"\n\
else\n\
    # Copy X11 authentication if available\n\
    if [ -f /tmp/.Xauthority ]; then\n\
        cp /tmp/.Xauthority /home/vscode/.Xauthority\n\
        chown vscode:vscode /home/vscode/.Xauthority\n\
    fi\n\
    exec "$@"\n\
fi' > /usr/local/bin/entrypoint.sh \
    && chmod +x /usr/local/bin/entrypoint.sh

# Set display for X11 forwarding
ENV DISPLAY=${DISPLAY:-:0}
ENV XAUTHORITY=/tmp/.Xauthority

# Switch back to vscode user
USER vscode
WORKDIR /workspaces/gimbal

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["sleep", "infinity"]
