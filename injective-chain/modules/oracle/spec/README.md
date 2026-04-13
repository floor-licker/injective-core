# `Oracle`

## Abstract

This specification describes the oracle module, which is primarily used by the `exchange` module to obtain external price data.

## Supported Oracle Types

The module currently supports the following active oracle types:

- **PriceFeed** — permissioned relayers submit prices directly for a base/quote pair.
- **Coinbase** — anyone can relay Coinbase-signed price messages (ECDSA verified against the Coinbase oracle address).
- **Provider** — permissioned relayers submit per-symbol prices under a named provider string.
- **Pyth** — the Pyth contract (configured in params) relays price attestations.
- **Stork** — permissioned publishers relay ECDSA-signed price messages.
- **Chainlink Data Streams** — any sender can submit Chainlink reports; each report is verified against the on-chain verifier proxy contract configured in params.

> **Note**: Band and Band IBC oracle types are deprecated and are no longer supported.

## Workflow

1. New oracle providers must first be authorized through a governance proposal that grants privileges to a list of relayers or publishers.
   - **PriceFeed**: `GrantPriceFeederPrivilegeProposal`
   - **Provider**: `GrantProviderPrivilegeProposal`
   - **Stork**: `GrantStorkPublisherPrivilegeProposal`
   - **Pyth** uses the module params (`MsgUpdateParams`) to configure the Pyth contract address; only that address can relay prices.
   - **Chainlink Data Streams** uses the module params to configure the verifier proxy contract address; any sender can submit reports, which are verified on-chain.
   - **Coinbase**: no proposal required — anyone can submit messages since they are signed by the Coinbase oracle key.

2. Once authorized, relayers submit oracle data using relay messages specific to their oracle type:
   - `MsgRelayPriceFeedPrice`
   - `MsgRelayProviderPrices`
   - `MsgRelayStorkPrices`
   - `MsgRelayPythPrices`
   - `MsgRelayCoinbaseMessages`
   - `MsgRelayChainlinkPrices`

3. Upon receiving a relay message, the oracle module verifies that the sender is authorized (or performs signature/contract verification), then persists the latest price data in state.

4. Other Cosmos SDK modules fetch the latest price data by querying the oracle module keeper (e.g. via `ViewKeeper.GetPrice`).

5. Module parameters (Pyth contract address, Chainlink verifier proxy contract address, and Chainlink verification gas limit) can be updated by the module authority via `MsgUpdateParams`.

**Note**: Privileges can be revoked through governance:
- `RevokePriceFeederPrivilegeProposal`
- `RevokeProviderPrivilegeProposal`
- `RevokeStorkPublisherPrivilegeProposal`

## Contents

1. [State](./01_state.md)
2. [Keeper](./02_keeper.md)
3. [Messages](./03_messages.md)
4. [Proposals](./04_proposals.md)
5. [Events](./05_events.md)
6. [Improvements](./06_future_improvements.md)
7. [Errors](./99_errors.md)
