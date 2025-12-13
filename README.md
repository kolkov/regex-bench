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
| literal_alt | 421 ms | 34 ms | **6 ms** | **12x faster** |
| anchored | 0.15 ms | **0.04 ms** | 0.31 ms | **4x faster** |
| inner_literal | 215 ms | 2.3 ms | **0.7 ms** | **94x faster** |
| suffix | 182 ms | 2.1 ms | **1.3 ms** | **87x faster** |
| char_class | 580 ms | **29 ms** | 65 ms | **20x faster** |
| email | 221 ms | 2.2 ms | **1.6 ms** | **99x faster** |

### Key Findings

**Go coregex v0.8.22 vs Go stdlib:**
- All patterns: **12-99x faster**
- Best: `email` **99x**, `inner_literal` **94x**, `suffix` **87x**

**Go coregex vs Rust regex:**
- `char_class`: **coregex 2.2x faster** (29ms vs 65ms)
- `anchored`: **coregex 7.8x faster** (0.04ms vs 0.31ms)
- `suffix`: Rust 1.6x faster
- `email`: Rust 1.4x faster
- `inner_literal`: Rust 3.3x faster
- `literal_alt`: Rust 5.7x faster (Aho-Corasick)

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
