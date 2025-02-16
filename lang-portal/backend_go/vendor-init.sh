#!/bin/bash
cd "$(dirname "$0")"
go mod vendor
go mod tidy