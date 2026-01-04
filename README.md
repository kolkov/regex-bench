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
| literal_alt | 475 ms | 32 ms | **0.9 ms** | **15x faster** |
| anchored | 0.02 ms | 0.28 ms | **0.09 ms** | — |
| inner_literal | 232 ms | **0.94 ms** | 0.74 ms | **247x faster** |
| suffix | 234 ms | 2.1 ms | **1.6 ms** | **112x faster** |
| char_class | 542 ms | 171 ms | **52 ms** | **3.2x faster** |
| email | 260 ms | 2.6 ms | **1.7 ms** | **101x faster** |
| uri | 257 ms | 1.9 ms | **1.1 ms** | **139x faster** |
| ip | 492 ms | 219 ms | **12 ms** | **2.2x faster** |

### Key Findings

**Go coregex v0.9.0 vs Go stdlib:**
- All patterns: **2-247x faster**
- Best: `inner_literal` **247x**, `uri` **139x**, `suffix` **112x**, `email` **101x**

**Go coregex vs Rust regex:**
- `inner_literal`: **coregex ~tie** (0.94ms vs 0.74ms)
- `suffix`: **coregex ~tie** (2.1ms vs 1.6ms)
- `email`: **coregex ~tie** (2.6ms vs 1.7ms)
- `uri`: Rust 1.7x faster
- `char_class`: Rust 3.3x faster
- `ip`: Rust 18x faster (complex alternation groups)
- `literal_alt`: Rust 36x faster (Aho-Corasick for 4 patterns)

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 2-247x slower |
| **Go coregex** | Reverse search, SIMD prefilters, Aho-Corasick (>8 patterns) | Slow on complex nested alternations |
| **Rust regex** | Aho-Corasick (any count), mature DFA | — |

**Note**: coregex v0.9.0 includes Aho-Corasick for >8 literal patterns and DigitPrefilter for digit-start patterns. These benchmarks use patterns that don't trigger these optimizations (literal_alt has only 4 patterns, ip uses complex nested alternations).

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
