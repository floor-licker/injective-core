---
sidebar_position: 4
title: Governance Proposals
---

# Governance Proposals

## GrantBandOraclePrivilegeProposal (Deprecated)

> **Deprecated.** Band oracle is no longer supported.

Band Oracle privileges could be granted to relayer accounts through a `GrantBandOraclePrivilegeProposal`.

```protobuf
message GrantBandOraclePrivilegeProposal {
    option (gogoproto.equal) = false;
    option (gogoproto.goproto_getters) = false;
    string title = 1;
    string description = 2;
    repeated string relayers = 3;
}
```

## RevokeBandOraclePrivilegeProposal (Deprecated)

> **Deprecated.** Band oracle is no longer supported.

Band Oracle privileges could be revoked from relayer accounts through a `RevokeBandOraclePrivilegeProposal`.

```protobuf
message RevokeBandOraclePrivilegeProposal {
    option (gogoproto.equal) = false;
    option (gogoproto.goproto_getters) = false;
    string title = 1;
    string description = 2;
    repeated string relayers = 3;
}
```

## GrantPriceFeederPrivilegeProposal

Price feeder privileges for a given base/quote pair are issued to relayers through a `GrantPriceFeederPrivilegeProposal`.

```protobuf
message GrantPriceFeederPrivilegeProposal {
    option (amino.name) = "oracle/GrantPriceFeederPrivilegeProposal";
    option (gogoproto.equal) = false;
    option (gogoproto.goproto_getters) = false;
    option (cosmos_proto.implements_interface) = "cosmos.gov.v1beta1.Content";
    string title = 1;
    string description = 2;
    string base = 3;
    string quote = 4;
    repeated string relayers = 5;
}
```

## RevokePriceFeederPrivilegeProposal

Price feeder privileges can be revoked from relayer accounts through a `RevokePriceFeederPrivilegeProposal`.

```protobuf
message RevokePriceFeederPrivilegeProposal {
    option (amino.name) = "oracle/RevokePriceFeederPrivilegeProposal";
    option (gogoproto.equal) = false;
    option (gogoproto.goproto_getters) = false;
    option (cosmos_proto.implements_interface) = "cosmos.gov.v1beta1.Content";
    string title = 1;
    string description = 2;
    string base = 3;
    string quote = 4;
    repeated string relayers = 5;
}
```

## GrantProviderPrivilegeProposal

Provider oracle privileges for a given provider name are issued to relayers through a `GrantProviderPrivilegeProposal`.

```protobuf
message GrantProviderPrivilegeProposal {
    option (amino.name) = "oracle/GrantProviderPrivilegeProposal";
    option (gogoproto.equal) = false;
    option (gogoproto.goproto_getters) = false;
    option (cosmos_proto.implements_interface) = "cosmos.gov.v1beta1.Content";
    string title = 1;
    string description = 2;
    string provider = 3;
    repeated string relayers = 4;
}
```

When accepted, the handler stores a `ProviderInfo` associating the given provider name with the list of authorized relayer addresses.

## RevokeProviderPrivilegeProposal

Provider oracle privileges can be revoked from relayer accounts through a `RevokeProviderPrivilegeProposal`.

```protobuf
message RevokeProviderPrivilegeProposal {
    option (amino.name) = "oracle/RevokeProviderPrivilegeProposal";
    option (gogoproto.equal) = false;
    option (gogoproto.goproto_getters) = false;
    option (cosmos_proto.implements_interface) = "cosmos.gov.v1beta1.Content";
    string title = 1;
    string description = 2;
    string provider = 3;
    repeated string relayers = 5;
}
```

## AuthorizeBandOracleRequestProposal (Deprecated)

> **Deprecated.** Band IBC oracle is no longer supported.

This proposal added a band oracle request to the list of IBC requests to be sent periodically.

```protobuf
message AuthorizeBandOracleRequestProposal {
    option (gogoproto.equal) = false;
    option (gogoproto.goproto_getters) = false;
    string title = 1;
    string description = 2;
    BandOracleRequest request = 3 [(gogoproto.nullable) = false];
}
```

## UpdateBandOracleRequestProposal (Deprecated)

> **Deprecated.** Band IBC oracle is no longer supported.

Used to delete or update an existing band oracle request. When `delete_request_id` is non-zero, the request with that ID is deleted. Otherwise, the request identified by `update_oracle_request.request_id` is updated.

```protobuf
message UpdateBandOracleRequestProposal {
    option (gogoproto.equal) = false;
    option (gogoproto.goproto_getters) = false;
    string title = 1;
    string description = 2;
    uint64 delete_request_id = 3;
    BandOracleRequest update_oracle_request = 4;
}
```

## EnableBandIBCProposal (Deprecated)

> **Deprecated.** Band IBC oracle is no longer supported.

This proposal enabled the IBC connection between Band chain and Injective chain and updated `BandIBCParams`.

```protobuf
message EnableBandIBCProposal {
    option (gogoproto.equal) = false;
    option (gogoproto.goproto_getters) = false;
    string title = 1;
    string description = 2;
    BandIBCParams band_ibc_params = 3 [(gogoproto.nullable) = false];
}
```

Details of `BandIBCParams` can be found in **[State](./01_state.md)**.

## GrantStorkPublisherPrivilegeProposal

Stork publisher privileges can be granted through a `GrantStorkPublisherPrivilegeProposal`.

```protobuf
message GrantStorkPublisherPrivilegeProposal {
    option (amino.name) = "oracle/GrantStorkPublisherPrivilegeProposal";
    option (gogoproto.equal) = false;
    option (gogoproto.goproto_getters) = false;
    option (cosmos_proto.implements_interface) = "cosmos.gov.v1beta1.Content";
    string title = 1;
    string description = 2;
    repeated string stork_publishers = 3;
}
```

## RevokeStorkPublisherPrivilegeProposal

Stork publisher privileges can be revoked through a `RevokeStorkPublisherPrivilegeProposal`.

```protobuf
message RevokeStorkPublisherPrivilegeProposal {
    option (amino.name) = "oracle/RevokeStorkPublisherPrivilegeProposal";
    option (gogoproto.equal) = false;
    option (gogoproto.goproto_getters) = false;
    option (cosmos_proto.implements_interface) = "cosmos.gov.v1beta1.Content";
    string title = 1;
    string description = 2;
    repeated string stork_publishers = 3;
}
```
