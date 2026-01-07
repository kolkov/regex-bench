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

| Pattern | Go stdlib | Go coregex | Rust regex | vs stdlib | vs Rust |
|---------|-----------|------------|------------|-----------|---------|
| suffix | 233 ms | **0.98 ms** | 1.34 ms | **238x** | **1.4x faster** |
| inner_literal | 231 ms | 1.43 ms | **0.55 ms** | **161x** | 2.6x slower |
| email | 260 ms | 1.51 ms | **1.40 ms** | **172x** | 8% slower |
| uri | 257 ms | 1.90 ms | **0.87 ms** | **135x** | 2.2x slower |
| ip | 494 ms | **3.80 ms** | 12.28 ms | **130x** | **3.2x faster** |
| multi_literal | 1436 ms | 12.42 ms | **4.91 ms** | **115x** | 2.5x slower |
| literal_alt | 474 ms | 4.27 ms | **0.80 ms** | **111x** | 5.3x slower |
| version | 168 ms | 2.17 ms | **0.66 ms** | **77x** | 3.3x slower |
| char_class | 521 ms | **35.71 ms** | 52.39 ms | **15x** | **1.5x faster** |
| anchored | 0.02 ms | 0.03 ms | 0.04 ms | ~1x | ~1x |

> **coregex v0.10.2** — version pattern regression fixed (was 8.2ms in v0.10.1, now 2.2ms).

### Key Findings

**Go coregex v0.10.2 vs Go stdlib:**
- Most patterns: **15-238x faster**
- Best: `suffix` **238x**, `email` **172x**, `inner_literal` **161x**, `uri` **135x**
- `ip` **130x** (DigitPrefilter)
- `multi_literal` **115x** (Aho-Corasick)
- `literal_alt` **111x** (Teddy SIMD)
- `version` **77x** (DigitPrefilter) — fixed in v0.10.2
- `char_class` **15x** (CharClassSearcher)

**Go coregex faster than Rust (3 patterns):**
- `ip`: **coregex 3.2x faster** (3.8ms vs 12.3ms)
- `char_class`: **coregex 1.5x faster** (36ms vs 52ms)
- `suffix`: **coregex 1.4x faster** (0.98ms vs 1.34ms)

**Rust faster than coregex:**
- `literal_alt`: Rust 5.3x faster (Teddy Fat with more buckets)
- `version`: Rust 3.3x faster
- `inner_literal`: Rust 2.6x faster
- `multi_literal`: Rust 2.5x faster
- `uri`: Rust 2.2x faster
- `email`: Rust 8% faster (almost equal)

> **Note**: Rust regex has 10+ years of development. coregex optimizations are targeted, not universal.

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 15-238x slower |
| **Go coregex** | Reverse search, SIMD prefilters, Aho-Corasick, **3 patterns faster than Rust** | Teddy gap vs Rust |
| **Rust regex** | Aho-Corasick (any count), mature DFA, overall fastest | char_class, ip, suffix slower than coregex |

**v0.10.2 (Current):**
- **Version pattern fixed**: 8.2ms → 2.2ms (DigitPrefilter restored)
- **3 patterns faster than Rust**: ip (3.2x), char_class (1.5x), suffix (1.4x)
- Gap vs Rust narrowing on most patterns

**Historical Improvements:**
- v0.10.0: Fat Teddy AVX2 (33-64 patterns, 9+ GB/s)
- v0.9.5: Aho-Corasick integration, Teddy 32 patterns
- v0.9.4: CharClassSearcher, Teddy 2-byte fingerprint
- v0.9.2: DigitPrefilter for IP patterns (3.2x faster than Rust)

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
| version | `\d+\.\d+\.\d+` | Version numbers | DigitPrefilter |
| ip | `(?:(?:25[0-5]\|2[0-4][0-9]\|...)\.){3}...` | IPv4 validation | DigitPrefilter + LazyDFA |

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
