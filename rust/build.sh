#!/usr/bin/env bash
set -e

cargo wasm-bundle --release
cp target/wasm-bundle/release/ising.html ./
