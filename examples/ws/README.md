# WebSocket API Examples

This directory contains examples of using the WebSocket API.

## Public API Example

The public API example shows how to use the WebSocket API to receive real-time market data:

- Book ticker updates

To run the example:

```bash
cd examples/ws
go run main.go
```

## Private API Example

The private API example shows how to use the WebSocket API to receive real-time account data:

- Account updates
- Position updates
- Order updates

To run the example:

```bash
cd examples/ws
go run main.go
```

Note: You need to set your API key and secret key in the example before running it.

## Testnet

All examples use the testnet by default. To use the mainnet, set `ws.UseTestnet = false` in the example.

The testnet endpoints are:
- WebSocket: `wss://testnet.blofin.com/ws`

The mainnet endpoints are:
- WebSocket: `wss://ws.blofin.com/ws` 