#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

ROOT_PACKAGE="github.com/lqshow/k8s-custom-controllers/foobar-code-generator"

CUSTOM_RESOURCE_NAME="foobar"

CUSTOM_RESOURCE_VERSION="v1alpha1"


"${GOPATH}"/src/k8s.io/code-generator/generate-groups.sh all \
    "${ROOT_PACKAGE}/pkg/generated" "${ROOT_PACKAGE}/pkg/apis" \
    "${CUSTOM_RESOURCE_NAME}:${CUSTOM_RESOURCE_VERSION}" \
    --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt
