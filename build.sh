#!/usr/bin/env bash

set -x

GitCommitLog=`git log --pretty=oneline -n 1`
GitCommitLog=${GitCommitLog//\'/\"}
GitStatus=`git status -s`

LDFlags=" \
    -X 'github.com/jhdriver/UBChain/param.GitCommitLog=${GitCommitLog}' \
    -X 'github.com/jhdriver/UBChain/param.GitStatus=${GitStatus}' \
"

ROOT_DIR=`pwd`
CHAIN_DIR=`pwd`
WALLET_DIR=`pwd`"/cmd/wallet"
BOOT_DIR=`pwd`"/cmd/tools"

rm -rf build
mkdir build


cd ${CHAIN_DIR} && GOOS=linux GOARCH=amd64 go build -ldflags "$LDFlags" -o ${ROOT_DIR}/build/linux/UBChain/UBChain &&
cd ${CHAIN_DIR} && GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFlags" -o ${ROOT_DIR}/build/darwin/UBChain/UBChain &&
cd ${CHAIN_DIR} && GOOS=windows GOARCH=amd64 go build -ldflags "$LDFlags" -o ${ROOT_DIR}/build/windows/UBChain/UBChain.exe &&
cp ${CHAIN_DIR}/config.toml ${ROOT_DIR}/build/linux/UBChain/ &&
cp ${CHAIN_DIR}/config.toml ${ROOT_DIR}/build/darwin/UBChain/ &&
cp ${CHAIN_DIR}/config.toml ${ROOT_DIR}/build/windows/UBChain/ &&

cd ${WALLET_DIR} && GOOS=linux GOARCH=amd64 go build -ldflags "$LDFlags" -o ${ROOT_DIR}/build/linux/wallet/wallet &&
cd ${WALLET_DIR} && GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFlags" -o ${ROOT_DIR}/build/darwin/wallet/wallet &&
cd ${WALLET_DIR} && GOOS=windows GOARCH=amd64 go build -ldflags "$LDFlags" -o ${ROOT_DIR}/build/windows/wallet/wallet.exe &&
cp ${WALLET_DIR}/wallet.toml ${ROOT_DIR}/build/linux/wallet/ &&
cp ${WALLET_DIR}/wallet.toml ${ROOT_DIR}/build/darwin/wallet/ &&
cp ${WALLET_DIR}/wallet.toml ${ROOT_DIR}/build/windows/wallet/ &&

cd ${BOOT_DIR} && GOOS=linux GOARCH=amd64 go build -ldflags "$LDFlags" -o ${ROOT_DIR}/build/linux/boot/boot &&
cd ${BOOT_DIR} && GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFlags" -o ${ROOT_DIR}/build/darwin/boot/boot &&
cd ${BOOT_DIR} && GOOS=windows GOARCH=amd64 go build -ldflags "$LDFlags" -o ${ROOT_DIR}/build/windows/boot/boot.exe &&



Version=`${ROOT_DIR}/build/darwin/UBChain/UBChain --version`

cd ${ROOT_DIR} &&
zip -r build/${Version}-linux-amd64.zip ./build/linux &&
zip -r build/${Version}-darwin-amd64.zip ./build/darwin &&
zip -r build/${Version}-windows-amd64.zip ./build/windows &&


ls -lrt ${ROOT_DIR}/build/linux/UBChain &&
ls -lrt ${ROOT_DIR}/build/linux/wallet &&
ls -lrt ${ROOT_DIR}/build/linux/boot &&

ls -lrt ${ROOT_DIR}/build/darwin/UBChain &&
ls -lrt ${ROOT_DIR}/build/darwin/wallet &&
ls -lrt ${ROOT_DIR}/build/darwin/boot &&

ls -lrt ${ROOT_DIR}/build/windows/UBChain &&
ls -lrt ${ROOT_DIR}/build/windows/wallet &&
ls -lrt ${ROOT_DIR}/build/windows/boot &&
echo 'build done.'