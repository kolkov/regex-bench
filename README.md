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
| literal_alt | 448 ms | 5.5 ms | **0.9 ms** | **81x faster** |
| multi_literal | 1250 ms | 47 ms | **4.7 ms** | **27x faster** |
| anchored | 0.05 ms | 0.45 ms | 0.05 ms | — |
| inner_literal | 201 ms | **1.9 ms** | 0.6 ms | **106x faster** |
| suffix | 203 ms | **1.4 ms** | 1.3 ms | **149x faster** |
| char_class | 494 ms | 63 ms | **53 ms** | **7.8x faster** |
| email | 245 ms | 2.0 ms | **1.6 ms** | **122x faster** |
| uri | 238 ms | 2.3 ms | **1.0 ms** | **103x faster** |
| version | 153 ms | 2.2 ms | **0.7 ms** | **70x faster** |
| **ip** | 457 ms | **3.2 ms** | 11.4 ms | **143x faster** |

### Key Findings

**Go coregex v0.9.4 vs Go stdlib:**
- Most patterns: **7-149x faster**
- Best: `suffix` **149x**, `ip` **143x**, `email` **122x**, `inner_literal` **106x**
- `literal_alt` **81x** (Teddy 2-byte fingerprint)
- `version` **70x** (DigitPrefilter)
- `char_class` **7.8x** (streaming state machine)

**Go coregex vs Rust regex:**
- `ip`: **coregex 3.6x faster** (3.2ms vs 11.4ms)
- `suffix`: ~tie (1.4ms vs 1.3ms)
- `char_class`: Rust 1.2x faster (was 2.6x in v0.9.2)
- `email`: Rust 1.2x faster
- `uri`: Rust 2.3x faster
- `inner_literal`: Rust 3.2x faster
- `version`: Rust 3.1x faster
- `literal_alt`: Rust 6.1x faster (was 39x in v0.9.2)
- `multi_literal`: Rust 10x faster (optimized Aho-Corasick)

> **Note**: Rust regex has 10+ years of development. coregex optimizations are targeted, not universal.

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 4-167x slower |
| **Go coregex** | Reverse search, SIMD prefilters, Aho-Corasick, **IP patterns faster than Rust** | Complex alternations |
| **Rust regex** | Aho-Corasick (any count), mature DFA, overall fastest | IP patterns slower than coregex |

**v0.9.4 Improvements:**
- 🚀 `char_class`: **7.8x faster** than stdlib (was 3.7x in v0.9.2)
  - Streaming state machine eliminates per-match overhead
  - Gap vs Rust reduced from 2.6x to 1.2x
- 🚀 `literal_alt`: **81x faster** than stdlib (was 15x in v0.9.2)
  - Teddy 2-byte fingerprint reduces false positives 90%
- 🚀 `version`: **70x faster** than stdlib (was 21x in v0.9.2)
  - DigitPrefilter prioritization over ReverseInner

**v0.9.2 Features:**
- ✅ `ip`: **143x faster** than stdlib, **3.6x faster than Rust**
- ✅ `multi_literal`: Aho-Corasick for >8 literal patterns
- ✅ DigitPrefilter for simple digit-lead patterns

## Patterns Tested

| Name | Pattern | Type | Optimization |
|------|---------|------|--------------|
| literal_alt | `error\|warning\|fatal\|critical` | 4-literal alternation | Teddy SIMD |
| multi_literal | `apple\|banana\|...\|orange` | 12-literal alternation | **Aho-Corasick** |
| anchored | `^HTTP/[12]\.[01]` | Start anchor | — |
| inner_literal | `.*@example\.com` | Inner literal | Reverse search |
| suffix | `.*\.(txt\|log\|md)` | Suffix match | Reverse search |
| char_class | `[\w]+` | Character class | CharClassSearcher |
| email | `[\w.+-]+@[\w.-]+\.[\w.-]+` | Complex real-world | Memmem SIMD |
| uri | `[\w]+://[^/\s?#]+[^\s?#]+...` | URL with query/fragment | Memmem SIMD |
| version | `\d+\.\d+\.\d+` | Version numbers | **ReverseInner** (`.` literal) |
| ip | `(?:(?:25[0-5]\|2[0-4][0-9]\|...)\.){3}...` | IPv4 validation | **LazyDFA** (optimized) |

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
