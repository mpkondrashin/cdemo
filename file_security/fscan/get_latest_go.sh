#!/bin/sh
curl -fsSL https://go.dev/dl/ | grep -oP 'go[0-9]+\.[0-9]+\.[0-9]+\.linux-amd64\.tar\.gz' | head -n1 | xargs -I {} sh -c '
  cd /tmp &&
  curl -LO https://go.dev/dl/{} &&
  sudo rm -rf /usr/local/go &&
  sudo tar -C /usr/local -xzf {} &&
  echo "export PATH=\$PATH:/usr/local/go/bin" | sudo tee /etc/profile.d/go.sh > /dev/null
'
