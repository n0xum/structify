#!/usr/bin/env bash
set -euo pipefail
: "${SSH_KEY:?SSH private key is required}" "${SSH_HOST:?SSH host is required}" "${SSH_PORT:?SSH port is required}" "${SSH_KNOWN_HOSTS:?pinned known_hosts is required}"
SSH_USER=${SSH_USER:-deploy}
[[ "$SSH_USER" == deploy ]]
[[ "$SSH_HOST" =~ ^[A-Za-z0-9.-]+$ && "$SSH_PORT" =~ ^[0-9]{1,5}$ ]]
(( SSH_PORT >= 1 && SSH_PORT <= 65535 ))
install -d -m 700 "$RUNNER_TEMP/ssh"
umask 077
printf '%s\n' "$SSH_KEY" > "$RUNNER_TEMP/ssh/id_key"
printf '%s\n' "$SSH_KNOWN_HOSTS" > "$RUNNER_TEMP/ssh/known_hosts"
cat > "$RUNNER_TEMP/ssh/config" <<EOF
Host gate
  HostName $SSH_HOST
  Port $SSH_PORT
  User $SSH_USER
  IdentityFile $RUNNER_TEMP/ssh/id_key
  IdentitiesOnly yes
  BatchMode yes
  StrictHostKeyChecking yes
  UserKnownHostsFile $RUNNER_TEMP/ssh/known_hosts
  GlobalKnownHostsFile /dev/null
  ForwardAgent no
  ForwardX11 no
  ClearAllForwardings yes
  RequestTTY no
  ConnectTimeout 15
EOF
chmod 600 "$RUNNER_TEMP/ssh/id_key" "$RUNNER_TEMP/ssh/known_hosts" "$RUNNER_TEMP/ssh/config"
