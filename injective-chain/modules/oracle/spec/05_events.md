---
sidebar_position: 5
title: Events
---

# Events

The oracle module emits the following events.

## Band (Deprecated)

> **Deprecated.** Band oracle events are no longer emitted for new price updates.

```protobuf
// Deprecated: Do not use.
message SetBandPriceEvent {
  string relayer = 1;
  string symbol = 2;
  string price = 3 [(gogoproto.customtype) = "cosmossdk.io/math.LegacyDec", (gogoproto.nullable) = false];
  uint64 resolve_time = 4;
  uint64 request_id = 5;
}

// Deprecated: Do not use.
message SetBandIBCPriceEvent {
  string relayer = 1;
  repeated string symbols = 2;
  repeated string prices = 3 [(gogoproto.customtype) = "cosmossdk.io/math.LegacyDec", (gogoproto.nullable) = false];
  uint64 resolve_time = 4;
  uint64 request_id = 5;
  int64 client_id = 6;
}

// Deprecated: Do not use.
message EventBandIBCAckSuccess {
  string ack_result = 1;
  int64 client_id = 2;
}

// Deprecated: Do not use.
message EventBandIBCAckError {
  string ack_error = 1;
  int64 client_id = 2;
}

// Deprecated: Do not use.
message EventBandIBCResponseTimeout {
  int64 client_id = 1;
}
```

## Legacy Chainlink (Deprecated)

> **Deprecated.** This event type has been superseded by `EventSetChainlinkDataStreamsPrices`.

```protobuf
// Deprecated: Do not use.
message SetChainlinkPriceEvent {
  string feed_id = 1;
  string answer = 2 [(gogoproto.customtype) = "cosmossdk.io/math.LegacyDec", (gogoproto.nullable) = false];
  uint64 timestamp = 3;
}
```

## Coinbase

```protobuf
message SetCoinbasePriceEvent {
  string symbol = 1;
  string price = 2 [(gogoproto.customtype) = "cosmossdk.io/math.LegacyDec", (gogoproto.nullable) = false];
  uint64 timestamp = 3;
}
```

## Provider

```protobuf
message SetProviderPriceEvent {
  string provider = 1;
  string relayer = 2;
  string symbol = 3;
  string price = 4 [(gogoproto.customtype) = "cosmossdk.io/math.LegacyDec", (gogoproto.nullable) = false];
}
```

## Pricefeed

```protobuf
message SetPriceFeedPriceEvent {
  string relayer = 1;
  string base = 2;
  string quote = 3;
  // price defines the price of the oracle base and quote
  string price = 4 [(gogoproto.customtype) = "cosmossdk.io/math.LegacyDec", (gogoproto.nullable) = false];
}
```

## Pyth

```protobuf
message EventSetPythPrices {
  repeated PythPriceState prices = 1;
}
```

## Stork

```protobuf
message EventSetStorkPrices {
  repeated StorkPriceState prices = 1;
}
```

## Chainlink Data Streams

```protobuf
message EventSetChainlinkDataStreamsPrices {
  repeated ChainlinkDataStreamsPriceState prices = 1;
}
```
