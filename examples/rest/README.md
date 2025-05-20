# REST API Examples

This directory contains examples of using the REST API.

## Public API Example

The public API example shows how to use the REST API to get market data:

- Funding rate history
- Candles (OHLCV data)
- Mark price

To run the example:

```bash
cd examples/rest
go run main.go
```

## Private API Example

The private API example shows how to use the REST API to manage orders:

- Account balance
- Leverage information
- Place order
- Get pending orders
- Cancel order

To run the example:

```bash
cd examples/rest
go run main.go
```

Note: You need to set your API key and secret key in the example before running it.

## Testnet

All examples use the testnet by default. To use the mainnet, set `rest.UseTestnet = false` in the example.

The testnet endpoints are:
- REST API: `https://testnet.blofin.com`

The mainnet endpoints are:
- REST API: `https://api.blofin.com` 