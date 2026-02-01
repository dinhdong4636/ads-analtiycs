#!/bin/bash
set -e

USER_ID=${UID:-1000}
GROUP_ID=${GID:-1000}
USERNAME="hostuser"
GROUPNAME="hostgroup"

# Create group if doesn't exist
if ! getent group ${GROUP_ID} >/dev/null 2>&1; then
    groupadd -g ${GROUP_ID} ${GROUPNAME} || true
fi

# Create user if doesn't exist  
if ! getent passwd ${USER_ID} >/dev/null 2>&1; then
    useradd -u ${USER_ID} -g ${GROUP_ID} -m -s /bin/bash -d /home/${USERNAME} ${USERNAME} || true
fi

# Install Air for the user (with proper PATH)
if [ ! -f "/go/bin/air" ]; then
    su - ${USERNAME} -c "PATH=/usr/local/go/bin:\$PATH GOPATH=/go go install github.com/cosmtrek/air@v1.49.0" || true
fi

# Symlink to global path
if [ -f "/go/bin/air" ]; then
    ln -sf /go/bin/air /usr/local/bin/air
fi

# Create and fix permissions
mkdir -p /app/export /app/bin /home/${USERNAME}
chown -R ${USER_ID}:${GROUP_ID} /app/export /app/bin /go /home/${USERNAME} 2>/dev/null || true
chown ${USER_ID}:${GROUP_ID} /app/go.mod /app/go.sum 2>/dev/null || true

echo "✓ Running as: ${USERNAME} (UID=${USER_ID}, GID=${GROUP_ID})"

# Execute command as user with proper environment
exec gosu ${USERNAME} env PATH=/usr/local/go/bin:/go/bin:$PATH GOPATH=/go HOME=/home/${USERNAME} "$@"
