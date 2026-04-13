---
sidebar_position: 3
title: Messages
---

# Messages

## MsgRelayBandRates (Deprecated)

> **Deprecated.** Band oracle is no longer supported.

Authorized Band relayers could relay price feed data for multiple symbols with `MsgRelayBandRates`. The handler iterated over all symbols and created/updated the `BandPriceState` for each.

```protobuf
message MsgRelayBandRates {
  string relayer = 1;
  repeated string symbols = 2;
  repeated uint64 rates = 3;
  repeated uint64 resolve_times = 4;
  repeated uint64 requestIDs = 5;
}
```

## MsgRelayCoinbaseMessages

Relayers of the Coinbase oracle can submit price data using `MsgRelayCoinbaseMessages`.

Each Coinbase message is authenticated by the `Signatures` provided by the Coinbase oracle address `0xfCEAdAFab14d46e20144F48824d0C09B1a03F2BC`, so anyone can submit this message.

```protobuf
message MsgRelayCoinbaseMessages {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  string sender = 1;
  repeated bytes messages = 2;
  repeated bytes signatures = 3;
}
```

This message fails if:
- Signature verification fails for any message.
- The timestamp submitted is strictly older than the last stored timestamp for that symbol. A message with a timestamp equal to the last stored timestamp is accepted as a no-op.

## MsgRelayPriceFeedPrice

Relayers of a PriceFeed oracle can relay prices using `MsgRelayPriceFeedPrice`.

```protobuf
message MsgRelayPriceFeedPrice {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  string sender = 1;
  repeated string base = 2;
  repeated string quote = 3;
  // price defines the price of the oracle base and quote
  repeated string price = 4 [
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable) = false
  ];
}
```

This message fails if:
- The sender is not an authorized PriceFeed relayer for the given base/quote pair.
- Any price is not positive or exceeds 10,000,000.

## MsgRequestBandIBCRates (Deprecated)

> **Deprecated.** Band IBC oracle is no longer supported.

`MsgRequestBandIBCRates` was used to instantly broadcast a price request to Band chain via IBC.

```protobuf
message MsgRequestBandIBCRates {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  string sender = 1;
  uint64 request_id = 2;
}
```

## MsgRelayPythPrices

`MsgRelayPythPrices` is sent by the Pyth contract to relay price attestations to the oracle module.

```protobuf
message MsgRelayPythPrices {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  string sender = 1;
  repeated PriceAttestation price_attestations = 2;
}

message PriceAttestation {
  string price_id = 1;
  int64 price = 2;
  uint64 conf = 3;
  int32 expo = 4;
  int64 ema_price = 5;
  uint64 ema_conf = 6;
  int32 ema_expo = 7;
  int64 publish_time = 8;
}
```

This message is a no-op (returns success) if the attestations list is empty.

This message fails if:
- The Pyth contract address is not configured in oracle module params.
- The sender does not equal the Pyth contract address defined in oracle module params.

Attestations that fail `PriceAttestation.Validate()` (e.g. bad price ID, exponent outside `[-12, 10]`) are skipped and logged as errors. Attestations with a publish time strictly older than the currently stored publish time, or whose price would move more than 100× from the last stored price, are silently skipped with no log. The message still succeeds as long as the sender check passes.

## MsgRelayStorkPrices

`MsgRelayStorkPrices` relays signed price messages from the Stork API to the oracle module.

```protobuf
message MsgRelayStorkPrices {
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1;
  repeated AssetPair asset_pairs = 2;
}

message AssetPair {
  string asset_id = 1;
  repeated SignedPriceOfAssetPair signed_prices = 2;
}

message SignedPriceOfAssetPair {
  string publisher_key = 1;
  uint64 timestamp = 2;
  string price = 3 [
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable) = false
  ];
  bytes signature = 4;
}
```

This message fails if:
- The `asset_id` values are not unique among the provided asset pairs.
- ECDSA signature verification fails for any `SignedPriceOfAssetPair` (verified in `ValidateBasic`).
- The difference between any two signed price timestamps exceeds `MaxStorkTimestampIntervalNano` (500 milliseconds).

Any sender may submit this message. At processing time, signed prices whose `publisher_key` is not in the authorized publisher list are silently skipped; the message succeeds as long as `ValidateBasic` passes.

## MsgRelayProviderPrices

Relayers of a provider-based oracle can submit prices using `MsgRelayProviderPrices`.

```protobuf
message MsgRelayProviderPrices {
  option (amino.name) = "oracle/MsgRelayProviderPrices";
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1;
  string provider = 2;
  repeated string symbols = 3;
  repeated string prices = 4 [
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable) = false
  ];
}
```

This message fails if:
- The sender is not an authorized relayer for the given provider.
- Any price is negative or exceeds 10,000,000 (zero prices are allowed for provider oracles).

## MsgRelayChainlinkPrices

`MsgRelayChainlinkPrices` relays Chainlink Data Streams reports to the oracle module.

```protobuf
message MsgRelayChainlinkPrices {
  option (amino.name) = "oracle/MsgRelayChainlinkPrices";
  option (gogoproto.equal) = false;
  option (gogoproto.goproto_getters) = false;
  option (cosmos.msg.v1.signer) = "sender";
  string sender = 1;
  repeated ChainlinkReport reports = 2;
}

message ChainlinkReport {
  bytes feed_id = 1;
  bytes full_report = 2;
  uint64 valid_from_timestamp = 3;
  uint64 observations_timestamp = 4;
}
```

Each report's `full_report` is verified on-chain using the Chainlink verifier proxy contract (`chainlink_verifier_proxy_contract` in params) via an EVM call. The EVM gas limit for verification is configurable via `chainlink_data_streams_verification_gas_limit` in params.

This message is a no-op (returns success) if the reports list is empty.

Each report is processed independently. This message fails only if every report fails processing (for any reason: report decoding, feed ID mismatch, or on-chain verification). If at least one report is processed successfully, the message succeeds regardless of how many others failed.

## MsgUpdateParams

`MsgUpdateParams` updates the oracle module parameters via governance authority.

```protobuf
message MsgUpdateParams {
  option (amino.name) = "oracle/MsgUpdateParams";
  option (cosmos.msg.v1.signer) = "authority";
  string authority = 1;
  Params params = 2 [(gogoproto.nullable) = false];
}
```

This message fails if:
- The sender is not the module authority (governance address).
- The params fail validation (e.g. invalid contract address format).

The params that can be updated are: `pyth_contract`, `chainlink_verifier_proxy_contract`, and `chainlink_data_streams_verification_gas_limit`.
