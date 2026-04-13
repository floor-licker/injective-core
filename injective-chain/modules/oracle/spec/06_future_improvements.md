---
sidebar_position: 6
title: Future Improvements
---

# Future Improvements

The oracle module currently supports the following oracle types:

- **PriceFeed** — permissioned relayers submit prices directly for a base/quote pair.
- **Coinbase** — anyone can relay Coinbase-signed price messages (ECDSA verified).
- **Provider** — permissioned relayers submit per-symbol prices under a named provider.
- **Pyth** — the Pyth contract relays price attestations on-chain.
- **Stork** — permissioned publishers relay ECDSA-signed price messages.
- **Chainlink Data Streams** — any sender can submit Chainlink reports; each report is verified on-chain against the verifier proxy contract configured in params.

The following oracle types are defined in the `OracleType` enum but are not yet integrated:

- **Razor** — not implemented.
- **DIA** — not implemented.
- **API3** — not implemented.
- **UMA** — not implemented.

The following oracle types were previously supported but are now deprecated:

- **Band** — direct Band relayer integration; deprecated and removed.
- **BandIBC** — Band oracle data via IBC; deprecated and removed.
- **Chainlink** (legacy) — the original Chainlink OCR integration; deprecated and superseded by Chainlink Data Streams.
