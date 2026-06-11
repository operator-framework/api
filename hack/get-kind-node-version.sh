#!/usr/bin/env bash

set -euo pipefail

# Determines the latest available kindest/node image tag for the Kubernetes
# minor version that matches the k8s.io/client-go dependency in go.mod.
#
# client-go v0.36.1 → look for kindest/node:v1.36.* → print "1.36.1"

MINOR=$(go list -m k8s.io/client-go | cut -d" " -f2 | sed 's/^v0\.\([0-9]\{1,\}\)\.[0-9]\{1,\}$/\1/')

if [[ -z "$MINOR" ]]; then
  echo "error: could not determine minor version from client-go" >&2
  exit 1
fi

TAG_PREFIX="v1.${MINOR}."

VERSION=$(
  curl -sf "https://hub.docker.com/v2/repositories/kindest/node/tags?page_size=100&name=${TAG_PREFIX}" |
    python3 -c "
import sys, json

data = json.load(sys.stdin)
tags = [r['name'] for r in data.get('results', []) if r['name'].startswith('${TAG_PREFIX}')]
if not tags:
    sys.exit(1)
tags.sort(key=lambda t: int(t.rsplit('.', 1)[-1]))
print(tags[-1].lstrip('v'))
"
)

if [[ -z "$VERSION" ]]; then
  echo "error: no kindest/node image found for v1.${MINOR}.*" >&2
  exit 1
fi

echo "$VERSION"
