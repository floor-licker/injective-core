#!/bin/sh -e

ABIGEN_VERSION=v1.16.3
abigen="go run github.com/ethereum/go-ethereum/cmd/abigen@$ABIGEN_VERSION"

solc --combined-json abi,bin \
  Bank.sol \
  CosmosTypes.sol \
  ReentrancyHook.sol \
  > ReentrancyHook.json
${abigen} --combined-json ReentrancyHook.json --pkg contracts --type ReentrancyHook --out ReentrancyHook.go
rm ReentrancyHook.json

exit 0
