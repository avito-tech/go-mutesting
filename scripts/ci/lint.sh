#!/bin/sh

if [ -z ${PKG+x} ]; then echo "PKG is not set"; exit 1; fi
if [ -z ${ROOT_DIR+x} ]; then echo "ROOT_DIR is not set"; exit 1; fi

# SA1019 (deprecated API) is suppressed until the go/loader -> go/packages migration is done.
echo "staticcheck:"
OUT=$(staticcheck -checks=all,-SA1019 $PKG/... 2>&1 | grep --invert-match -E "(/example)")
if [ -n "$OUT" ]; then echo "$OUT"; PROBLEM=1; fi

if [ -n "$PROBLEM" ]; then exit 1; fi
