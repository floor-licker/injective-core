# Exchange Module Benchmark Suite

Comprehensive benchmarks for measuring exchange module performance under realistic production conditions.

## Rationale

Finding optimisation targets requires reproducible, realistic workloads that exercise the full production code path. This suite simulates realistic trading activity (market makers refreshing quotes, traders placing/cancelling orders, liquidations, funding) while executing the complete transaction lifecycle: `Tx Build → Sign (secp256k1) → ProcessProposal → FinalizeBlock → Commit`. The result is CPU profiles that (hopefully) accurately reflect production hotspots.

## How It Works

1. **Setup**: Creates configurable spot + derivative markets, funds trader/MM subaccounts, warms up orderbooks with resting liquidity
2. **Per-block simulation**:
   - Price engine advances (geometric Brownian motion with mean reversion)
   - Oracle prices update (with configurable divergence to trigger funding)
   - Market makers refresh quotes via `MsgBatchUpdateOrders` (cancel-all + place new)
   - Traders submit orders (limit/market/conditional/post-only/reduce-only) and cancellations
   - Liquidations attempted when oracle moves trigger underwater positions
3. **Block execution**: Pending messages are signed, batched into a block, and executed through `ProcessProposal → FinalizeBlock → Commit`
4. **Metrics collection**: Gas usage, throughput, order outcomes, and execution timings are recorded for analysis

The benchmark timer measures only block execution (step 3), excluding workload generation and result observation, ensuring metrics reflect actual chain performance.

## Quick Start

```bash
# Run from the benchmark package directory
cd injective-chain/modules/exchange/keeper/benchmarks

# Add -v to include detailed benchmark reports (b.Log output)

# CI (~5s)
go test -bench=BenchmarkRealisticTrading_Light -benchtime=1x

# Performance analysis (~2–5min, 50k orders)
go test -bench=BenchmarkRealisticTrading_Heavy -benchtime=1x -timeout=30m

# Extended profiling (~10–20min, 300k orders)
go test -bench=BenchmarkRealisticTrading_Extended -benchtime=1x -timeout=30m
```

## Configuration

### Presets

| Preset | Blocks | Orders/Block | Total Orders | Markets (spot+deriv) | Use Case |
|--------|--------|--------------|--------------|----------------------|----------|
| Light | 10 | 10 | 100 | 1+1 | CI |
| Heavy | 100 | 500 | 50,000 | 3+7 | Performance analysis |
| Extended | 1000 | 300 | 300,000 | 3+5 | Extended runs |
| Extreme | 500 | 5000 | 2,500,000 | 10+20 | Stress testing |

### JSON Config File

For full control, use a JSON config file:

```bash
BENCH_CONFIG=./config.json go test -bench=BenchmarkRealisticTrading_Custom -benchtime=1x
```

See `config.example.json` for the full schema.

### Environment Variable Overrides

Override individual parameters on any preset:

```bash
BENCH_BLOCKS=100 BENCH_ORDERS_PER_BLOCK=200 \
  go test -bench=BenchmarkRealisticTrading_Heavy -benchtime=1x
```

| Variable | Description |
|----------|-------------|
| `BENCH_CONFIG` | Path to JSON config file |
| `BENCH_PRESET` | Base preset (light/heavy/extended/extreme) |
| `BENCH_BLOCKS` | Block count |
| `BENCH_ORDERS_PER_BLOCK` | Orders per block |
| `BENCH_SPOT_MARKETS` | Spot market count (0 allowed) |
| `BENCH_DERIVATIVE_MARKETS` | Derivative market count (0 allowed) |
| `BENCH_TRADERS` | Trader count |
| `BENCH_MARKET_MAKERS` | Market maker count |
| `BENCH_ORDERBOOK_DEPTH` | Orders per side in warmup |
| `BENCH_SECONDS_PER_BLOCK` | Simulated seconds per block (for funding rate testing) |
| `BENCH_LIMIT_PCT` | Limit order percentage (0-100) |
| `BENCH_MARKET_PCT` | Market order percentage (0-100) |
| `BENCH_BATCH_ORDER_PCT` | Probability (0-100) that a block includes additional batch orders |
| `BENCH_BATCH_SIZE` | Total orders per batch message (0 = disabled, capped at 100) |
| `BENCH_EXPIRING_ORDER_PCT` | % of limit orders with expiration |
| `BENCH_TX_GAS_LIMIT` | Gas limit per transaction (0 = use default 500000) |
| `BENCH_SEED` | Random seed for reproducible runs (0 = non-deterministic) |
| `BENCH_LIQUIDATION_TRIGGER_PCT` | Probability (0-1) of liquidation-level price moves per block |
| `BENCH_ORACLE_DIVERGENCE_PCT` | Max oracle divergence from mid price (0-1, e.g. 0.03 = 3%) |
| `BENCH_ORACLE_TREND_BIAS` | Oracle trend direction bias (-1 to 1) |
| `BENCH_ENABLE_EXTREME` | Set to "1" to enable Extreme benchmark (prevents accidental OOM) |
| `BENCHMARK_REPORT_DIR` | Directory for output files (JSON/CSV/TXT reports) |
| `BENCHMARK_REPORT_PREFIX` | Custom prefix for report filenames (e.g. "heavy_seed42") |

## Reproducibility

For deterministic benchmark runs, set a specific seed:

```bash
BENCH_SEED=12345 go test -bench=BenchmarkRealisticTrading_Heavy -benchtime=1x
```

The seed used is logged at the start of each run:
```
Simulation seed: 12345 (reproduce with BENCH_SEED=12345)
```

To reproduce a previous run, use the seed from its output.

## Profiling

### Generate Profiles

```bash
go test -bench=BenchmarkRealisticTrading_Heavy -benchtime=1x \
  -cpuprofile=cpu.prof -memprofile=mem.prof -timeout=30m
```

### Analyse Profiles

```bash
# Interactive (opens browser)
go tool pprof -http=:8080 cpu.prof
go tool pprof -http=:8081 mem.prof

# Text output
go tool pprof -top cpu.prof
go tool pprof -top mem.prof

# Flamegraph
go tool pprof -svg cpu.prof > cpu.svg
```

## Output

### Reports

Reports are written to stdout. To save to files:

```bash
BENCHMARK_REPORT_DIR=./reports go test -bench=BenchmarkRealisticTrading_Heavy -benchtime=1x
```

Creates timestamped files (e.g., `20260131_143052_`):
- `reports/<timestamp>_benchmark_report.txt` - Human-readable report
- `reports/<timestamp>_benchmark_metrics.json` - Machine-readable metrics
- `reports/<timestamp>_benchmark_metrics.csv` - CSV export

Set `BENCHMARK_REPORT_PREFIX` to add a custom prefix (e.g., `heavy_seed42_`).

### Metrics

**Order metrics:**
- Orders: created, matched, cancelled, errored
- Batch orders created
- Batch cancellations (markets cancelled via CancelAll)
- Expiring orders created
- Orders expired

**Position metrics:**
- Conditional orders created and triggered
- Liquidations executed
- Funding payments applied

**Account operations:**
- Deposits, withdrawals, subaccount transfers

**Performance metrics:**
- Throughput (orders/sec, blocks/sec) based on block execution time (ProcessProposal + FinalizeBlock + Commit). Tx build/sign and result observation are excluded.
- Gas usage (total, per block, per order)
- Transaction success/failure rates

## Features Exercised

### Order Types
- Limit orders (70% default)
- Market orders (20% default)
- Conditional orders (10% default): STOP_BUY, STOP_SELL, TAKE_BUY, TAKE_SELL
- Post-only orders (~10% of limit orders)
- Reduce-only orders (~20% when trader has position)
- Orders with expiration (~5% of limit orders)

### Batch Operations
- `MsgBatchUpdateOrders` (cancel-all + create in single tx) - primary MM message type

Note: The realistic trading simulation primarily uses `MsgBatchUpdateOrders` with `SpotMarketIdsToCancelAll` / `DerivativeMarketIdsToCancelAll` for MM refresh.

### Account Operations
- `MsgDeposit` (~3% of blocks)
- `MsgWithdraw` (~2% of blocks)
- `MsgSubaccountTransfer` (~2% of blocks)

### Lifecycle Features
- Full block lifecycle (BeginBlocker → DeliverTx → EndBlocker → Commit)
- Transaction signature verification (secp256k1, not bypassed)
- IAVL state persistence
- Funding rate triggers via time acceleration (SecondsPerBlock=120)
- Liquidation attempts via `MsgLiquidatePosition` when oracle price moves significantly
- Order expiration processing
- Conditional order triggering

## Execution Models

### Realistic Trading Simulations (Recommended for Profiling)

The `BenchmarkRealisticTrading_*` benchmarks exercise the production code path:

1. **Transaction Building**: Messages are wrapped in signed transactions with proper sequence numbers
2. **ProcessProposal**: Transactions are validated
3. **FinalizeBlock**: Full block execution (PreBlocker → BeginBlocker → DeliverTx → EndBlocker)
4. **Commit**: State persisted to IAVL tree
