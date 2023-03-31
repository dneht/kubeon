/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package configset

const apiserverStartupService = `[Unit]
Description=Change host when worker startup
ConditionPathExists=/opt/kubeon/apiserver-startup.sh
Wants=network-online.target
After=network-online.target
Before=kubelet.target

[Service]
Type=forking
ExecStart=/opt/kubeon/apiserver-startup.sh
TimeoutSec=0
StandardOutput=tty
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
`

const apiserverStartupBash = `#!/usr/bin/env bash
TARGET_DOMAIN={{.TargetDomain}}
VIRTUAL_ADDR={{.VirtualAddr}}
REAL_ADDRS={{.RealAddrs}}
DEFAULT_PORT=6443

function check_server() {
  IN_ADDR=$1
  IN_PORT=$2
  if [[ -e "/usr/bin/curl" ]]; then
    curl -m 1 -s -k "https://${IN_ADDR}:${IN_PORT}/healthz" || echo -n no
  elif [[ -e "/usr/bin/nc" ]]; then
    nc -w 1 -z "${IN_ADDR}" "${IN_PORT}" && echo -n ok || echo -n no
  elif [[ -e "/usr/bin/timeout" ]]; then
    timeout 1 bash -c "</dev/tcp/${IN_ADDR}/${IN_PORT} &>/dev/null" && echo -n ok || echo -n no
  else
    bash -c "</dev/tcp/${IN_ADDR}/${IN_PORT} &>/dev/null" && echo -n ok || echo -n no
  fi
}

function update_host() {
    IN_ADDR=$1
    sed -i -E "s/^[0-9a-f.:]+\s+${TARGET_DOMAIN}.*$/${IN_ADDR}  ${TARGET_DOMAIN}/g" /etc/hosts
}

if [[ "$(check_server ${VIRTUAL_ADDR} ${DEFAULT_PORT})" == "ok" ]]; then
  echo "virtual ip can provide services normally"
  exit 0
fi

IS_ALREADY=""
for _ in {1..3}; do
  if [[ "${IS_ALREADY}" == "ok" ]]; then
    break
  fi

  for REAL_ADDR in ${REAL_ADDRS//,/ }; do
    if [[ "$(check_server ${REAL_ADDR} ${DEFAULT_PORT})" == "ok" ]]; then
      echo "replace host with available real ip ${REAL_ADDR}"
      update_host ${REAL_ADDR}
      IS_ALREADY="ok"
      break
    fi
  done
done

for _ in {1..300}; do
  sleep 5s
  if [[ "$(check_server ${VIRTUAL_ADDR} ${DEFAULT_PORT})" == "ok" ]]; then
    echo "virtual ip is available, task completed"
    update_host ${VIRTUAL_ADDR}
    exit 0
  else
    echo "virtual ip is not available, please wait..."
  fi
done
`
