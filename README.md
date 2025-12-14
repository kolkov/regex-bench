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

**GitHub Actions Ubuntu, 6.0 MB input** (using `FindAll` for fair comparison)

| Pattern | Go stdlib | Go coregex | Rust regex | coregex vs stdlib |
|---------|-----------|------------|------------|-------------------|
| literal_alt | 473 ms | 31 ms | **0.7 ms** | **15x faster** |
| anchored | 0.02 ms | 0.19 ms | **0.04 ms** | — |
| inner_literal | 231 ms | 1.9 ms | **0.6 ms** | **122x faster** |
| suffix | 233 ms | **1.8 ms** | 1.4 ms | **127x faster** |
| char_class | 525 ms | 119 ms | **52 ms** | **4.4x faster** |
| email | 259 ms | 1.7 ms | **1.3 ms** | **155x faster** |
| uri | 266 ms | 2.8 ms | **0.9 ms** | **96x faster** |
| ip | 493 ms | 164 ms | **12 ms** | **3x faster** |

### Key Findings

**Go coregex v0.8.24 vs Go stdlib:**
- All patterns: **3-155x faster**
- Best: `email` **155x**, `suffix` **127x**, `inner_literal` **122x**, `uri` **96x**

**Go coregex vs Rust regex:**
- `suffix`: **coregex ~tie** (1.8ms vs 1.4ms)
- `email`: **coregex ~tie** (1.7ms vs 1.3ms)
- `char_class`: Rust 2.3x faster
- `uri`: Rust 3x faster
- `inner_literal`: Rust 3x faster
- `ip`: Rust 14x faster
- `literal_alt`: Rust 44x faster (Aho-Corasick)

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 3-155x slower |
| **Go coregex** | Reverse search, SIMD prefilters | No Aho-Corasick, slow complex alternations |
| **Rust regex** | Aho-Corasick, mature optimizations | — |

Rust's advantage on `literal_alt` comes from Aho-Corasick multi-pattern matching. The `ip` pattern shows coregex weakness with complex alternation groups. coregex excels at suffix/inner literal patterns due to reverse search optimization.

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
