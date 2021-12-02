#!/bin/sh

if [ -z ${PKG+x} ]; then echo "PKG is not set"; exit 1; fi
if [ -z ${ROOT_DIR+x} ]; then echo "ROOT_DIR is not set"; exit 1; fi

echo "go vet:"
OUT=$(go vet -all=true ./... 2>&1 | grep --invert-match -E "(Checking file|\%p of wrong type|can't check non-constant format|/example)")
if [ -n "$OUT" ]; then echo "$OUT"; PROBLEM=1; fi

if [ -n "$PROBLEM" ]; then exit 1; fi
