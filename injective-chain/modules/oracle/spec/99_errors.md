# Error Codes

This document lists the error codes used in the oracle module.

> **Note on errors 36–38**: The registered error string for codes 36, 37, and 38 in the current implementation is `"unauthorized Pyth price relay"` (a copy-paste from code 35). The descriptions below reflect the **semantic intent** of each error variable (`ErrInvalidPythPriceID`, `ErrInvalidPythExponent`, `ErrInvalidPythPublishTime`), not the literal text that appears in emitted error messages.

| Module | Error Code | Description |
|--------|------------|-------------|
| oracle | 1  | relayer address is empty |
| oracle | 2  | bad rates count |
| oracle | 3  | bad resolve times |
| oracle | 4  | bad request ID |
| oracle | 5  | relayer not authorized |
| oracle | 6  | bad price feed base count |
| oracle | 7  | bad price feed quote count |
| oracle | 8  | unsupported oracle type |
| oracle | 9  | bad messages count |
| oracle | 10 | bad Coinbase message |
| oracle | 11 | bad Ethereum signature |
| oracle | 12 | bad Coinbase message timestamp |
| oracle | 13 | Coinbase price not found |
| oracle | 14 | prices must be positive |
| oracle | 15 | prices must be less than 10 million |
| oracle | 16 | invalid Band IBC request |
| oracle | 17 | sample error |
| oracle | 18 | invalid packet timeout |
| oracle | 19 | invalid symbols count |
| oracle | 20 | could not claim port capability |
| oracle | 21 | invalid IBC Port ID |
| oracle | 22 | invalid IBC Channel ID |
| oracle | 23 | invalid Band IBC request interval |
| oracle | 24 | invalid Band IBC update request proposal |
| oracle | 25 | Band IBC oracle request not found |
| oracle | 26 | base info is empty |
| oracle | 27 | provider is empty |
| oracle | 28 | invalid provider name |
| oracle | 29 | invalid symbol |
| oracle | 30 | relayer already exists |
| oracle | 31 | provider price not found |
| oracle | 32 | invalid oracle request |
| oracle | 33 | no price for oracle was found |
| oracle | 34 | no address for Pyth contract found |
| oracle | 35 | unauthorized Pyth price relay |
| oracle | 36 | invalid Pyth price ID |
| oracle | 37 | invalid Pyth exponent |
| oracle | 38 | invalid Pyth publish time |
| oracle | 39 | empty price attestations |
| oracle | 40 | bad Stork message timestamp |
| oracle | 41 | sender stork is empty |
| oracle | 42 | invalid stork signature |
| oracle | 43 | stork asset id not unique |
| oracle | 44 | chainlink report verification failed |
