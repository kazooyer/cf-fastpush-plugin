#!/bin/bash

set -e

(cf uninstall-plugin "cf-fastpush" || true) && go build -gcflags="-trimpath=${HOME}" -o cf-fastpush main.go && cf install-plugin cf-fastpush
