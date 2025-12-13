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

**GitHub Actions Ubuntu, 6.0 MB input**

| Pattern | Go stdlib | Go coregex | Rust regex | coregex vs stdlib |
|---------|-----------|------------|------------|-------------------|
| literal_alt | 473 ms | 31 ms | **0.8 ms** | **15x faster** |
| anchored | 0.02 ms | **0.01 ms** | 0.10 ms | **2x faster** |
| inner_literal | 232 ms | 1.5 ms | **0.6 ms** | **153x faster** |
| suffix | 240 ms | **1.5 ms** | 1.3 ms | **166x faster** |
| char_class | 550 ms | **26 ms** | 52 ms | **21x faster** |
| email | 259 ms | **1.5 ms** | 1.5 ms | **172x faster** |
| uri | 257 ms | 1.3 ms | **0.8 ms** | **192x faster** |
| ip | 493 ms | 163 ms | **12 ms** | **3x faster** |

### Key Findings

**Go coregex v0.8.22 vs Go stdlib:**
- All patterns: **3-192x faster**
- Best: `uri` **192x**, `email` **172x**, `suffix` **166x**, `inner_literal` **153x**

**Go coregex vs Rust regex:**
- `char_class`: **coregex 2x faster** (26ms vs 52ms)
- `anchored`: **coregex 10x faster** (0.01ms vs 0.10ms)
- `suffix`: **coregex 1.1x faster** (1.5ms vs 1.3ms - tie)
- `email`: **tie** (1.5ms vs 1.5ms)
- `uri`: Rust 1.6x faster
- `inner_literal`: Rust 2.4x faster
- `ip`: Rust 13x faster
- `literal_alt`: Rust 39x faster (Aho-Corasick)

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 3-192x slower |
| **Go coregex** | CharClassSearcher, reverse search | No Aho-Corasick, slow complex alternations |
| **Rust regex** | Aho-Corasick, mature optimizations | — |

Rust's advantage on `literal_alt` comes from Aho-Corasick multi-pattern matching. coregex wins on character classes due to CharClassSearcher's 256-byte lookup table. The `ip` pattern shows coregex weakness with complex alternation groups.

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
