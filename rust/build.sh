#!/usr/bin/env bash
set -e

PROFILE="release"
ARGS="--release"
if [[ $* == *--debug* ]]; then
    PROFILE="debug"
    ARGS=""
fi

cargo +nightly build --target wasm32-unknown-unknown ${ARGS}
wasm-bindgen --target web target/wasm32-unknown-unknown/${PROFILE}/ising.wasm --out-dir target/wasm32-unknown-unknown/${PROFILE}/
cargo build ${ARGS} --example pack
(
    cd target/wasm32-unknown-unknown/${PROFILE}/
    ../../${PROFILE}/examples/pack
)
cp target/wasm32-unknown-unknown/${PROFILE}/main.html ./