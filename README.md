# regex-bench

Cross-language regex benchmark for **real-world patterns**.

Created to provide data for [golang/go#26623](https://github.com/golang/go/issues/26623) discussion on Go regex performance.

## Test Environment

All benchmarks run on **identical conditions**:
- **OS**: Linux (Ubuntu via WSL2 or GitHub Actions)
- **Input**: 6.0 MB generated text file
- **Method**: Each engine compiled natively, same input file, same patterns

> **Note**: Cross-compiled Go binaries run in WSL2 for fair comparison with Rust.

## Results

**Intel i7-1255U, 6.0 MB input, Linux/WSL2**

| Pattern | Go stdlib | Go coregex | Rust regex | coregex vs stdlib |
|---------|-----------|------------|------------|-------------------|
| literal_alt | 410 ms | 36 ms | **6 ms** | **11x faster** |
| anchored | 0.15 ms | **0.03 ms** | 0.31 ms | **5x faster** |
| inner_literal | 220 ms | 2.0 ms | **0.7 ms** | **110x faster** |
| suffix | 185 ms | 1.3 ms | **1.3 ms** | **142x faster** |
| char_class | 466 ms | **24 ms** | 65 ms | **19x faster** |
| email | 206 ms | 1.3 ms | **1.6 ms** | **158x faster** |
| uri | 210 ms | 1.7 ms | *TBD* | **123x faster** |
| ip | 374 ms | 120 ms | *TBD* | **3x faster** |

### Key Findings

**Go coregex v0.8.22 vs Go stdlib:**
- All patterns: **3-158x faster**
- Best: `email` **158x**, `suffix` **142x**, `uri` **123x**, `inner_literal` **110x**

**Go coregex vs Rust regex:**
- `char_class`: **coregex 2.7x faster** (24ms vs 65ms)
- `anchored`: **coregex 10x faster** (0.03ms vs 0.31ms)
- `suffix`: **tie** (1.3ms vs 1.3ms)
- `email`: **coregex 1.2x faster** (1.3ms vs 1.6ms)
- `inner_literal`: Rust 2.9x faster
- `literal_alt`: Rust 6x faster (Aho-Corasick)

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 12-99x slower |
| **Go coregex** | CharClassSearcher, reverse search | No Aho-Corasick yet |
| **Rust regex** | Aho-Corasick, mature optimizations | — |

Rust's advantage on `literal_alt` comes from Aho-Corasick multi-pattern matching. coregex wins on character classes due to CharClassSearcher's 256-byte lookup table.

## Patterns Tested

| Name | Pattern | Type |
|------|---------|------|
| literal_alt | `error\|warning\|fatal\|critical` | Multi-literal alternation |
| anchored | `^HTTP/[12]\.[01]` | Start anchor |
| inner_literal | `.*@example\.com` | Inner literal (reverse search) |
| suffix | `.*\.(txt\|log\|md)` | Suffix match (reverse search) |
| char_class | `[\w]+` | Character class |
| email | `[\w.+-]+@[\w.-]+\.[\w.-]+` | Complex real-world |
| uri | `[\w]+://[^/\s?#]+[^\s?#]+...` | URL with query/fragment |
| ip | `(?:(?:25[0-5]\|2[0-4][0-9]\|...)\.){3}...` | IPv4 validation |

## Running Benchmarks

```bash
# Generate input data (6 MB)
go run scripts/generate-input.go

# Build for Linux
cd go-stdlib && GOOS=linux GOARCH=amd64 go build -o ../bin/go-stdlib-linux . && cd ..
cd go-coregex && GOOS=linux GOARCH=amd64 go build -o ../bin/go-coregex-linux . && cd ..

# Run all in WSL/Linux for fair comparison
wsl ./bin/go-stdlib-linux input/data.txt
wsl ./bin/go-coregex-linux input/data.txt
wsl ./bin/rust-benchmark input/data.txt
```

## CI Benchmarks

Benchmarks run automatically on GitHub Actions (Ubuntu) for reproducible results.

[![Benchmark](https://github.com/kolkov/regex-bench/actions/workflows/benchmark.yml/badge.svg)](https://github.com/kolkov/regex-bench/actions/workflows/benchmark.yml)

## Links

- **coregex**: https://github.com/coregx/coregex
- **Go issue**: https://github.com/golang/go/issues/26623
- **Rust regex**: https://github.com/rust-lang/regex

## License

MIT
