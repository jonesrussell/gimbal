ARG VARIANT=1.25-bookworm
FROM mcr.microsoft.com/devcontainers/go:${VARIANT}

# Remove expired Yarn repository
RUN rm -f /etc/apt/sources.list.d/yarn.list || true

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
    libpng-dev \
    libasound2-dev \
    alsa-utils \
    pulseaudio \
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
RUN printf '%s\n' \
    '#!/bin/bash' \
    'if [ "$1" = "test" ]; then' \
    '    # Set up fake X server with authentication' \
    '    export XAUTHORITY=/tmp/.Xauthority' \
    '    touch $XAUTHORITY' \
    '    xauth generate :99 . trusted' \
    '    Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &' \
    '    export DISPLAY=:99' \
    '    shift' \
    '    exec "$@"' \
    'else' \
    '    # Copy X11 authentication if available' \
    '    if [ -f /tmp/.Xauthority ]; then' \
    '        cp /tmp/.Xauthority /home/vscode/.Xauthority' \
    '        chown vscode:vscode /home/vscode/.Xauthority' \
    '    fi' \
    '    exec "$@"' \
    'fi' \
    > /usr/local/bin/entrypoint.sh \
    && chmod +x /usr/local/bin/entrypoint.sh

# Set display for X11 forwarding
ENV DISPLAY=${DISPLAY:-:0}
ENV XAUTHORITY=/tmp/.Xauthority

# Set up dummy ALSA device for audio (prevents errors in containers without audio)
# This creates a null audio device that won't fail when opened
RUN mkdir -p /home/vscode && \
    echo 'pcm.!default {' > /home/vscode/.asoundrc && \
    echo '  type plug' >> /home/vscode/.asoundrc && \
    echo '  slave.pcm "null"' >> /home/vscode/.asoundrc && \
    echo '}' >> /home/vscode/.asoundrc && \
    echo 'pcm.null {' >> /home/vscode/.asoundrc && \
    echo '  type file' >> /home/vscode/.asoundrc && \
    echo '  file /dev/null' >> /home/vscode/.asoundrc && \
    echo '}' >> /home/vscode/.asoundrc && \
    chown vscode:vscode /home/vscode/.asoundrc

# Switch back to vscode user
# USER vscode

WORKDIR /workspaces/gimbal
