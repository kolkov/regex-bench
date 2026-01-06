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
| literal_alt | 433 ms | 4.2 ms | **1.0 ms** | **104x** ✅ | 4.3x slower |
| multi_literal | 1269 ms | 12.9 ms | **4.7 ms** | **99x** ✅ | 2.8x slower |
| anchored | 0.05 ms | 0.53 ms | 0.06 ms | — | — |
| inner_literal | 203 ms | **1.5 ms** | 0.6 ms | **140x** ✅ | 2.3x slower |
| suffix | 203 ms | **1.3 ms** | 1.5 ms | **158x** ✅ | **~tie** |
| char_class | 506 ms | 56.5 ms | **51 ms** | **9x** ✅ | 1.1x slower |
| email | 244 ms | 1.7 ms | **1.5 ms** | **145x** ✅ | 1.1x slower |
| uri | 237 ms | 2.1 ms | **1.0 ms** | **111x** ✅ | 2.2x slower |
| version | 152 ms | 1.6 ms | **0.8 ms** | **93x** ✅ | 2.2x slower |
| **ip** | 458 ms | **2.9 ms** | 11.4 ms | **157x** ✅ | **3.9x faster** ✅ |

### Key Findings

**Go coregex v0.9.5 vs Go stdlib:**
- Most patterns: **9-158x faster**
- Best: `suffix` **158x**, `ip` **157x**, `email` **145x**, `inner_literal` **140x**
- `multi_literal` **99x** (Aho-Corasick + literal extractor fix)
- `literal_alt` **104x** (Teddy SIMD)
- `version` **93x** (DigitPrefilter)
- `char_class` **9x** (streaming state machine)

**Go coregex vs Rust regex:**
- `ip`: **coregex 3.9x faster** (2.9ms vs 11.4ms)
- `suffix`: ~tie
- `char_class`: Rust 1.1x faster
- `email`: Rust 1.1x faster
- `uri`: Rust 2.2x faster
- `inner_literal`: Rust 2.3x faster
- `version`: Rust 2.2x faster
- `literal_alt`: Rust 4.3x faster (was 6.1x in v0.9.4)
- `multi_literal`: Rust 2.8x faster (was 10x in v0.9.4!)

> **Note**: Rust regex has 10+ years of development. coregex optimizations are targeted, not universal.

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 9-158x slower |
| **Go coregex** | Reverse search, SIMD prefilters, Aho-Corasick, **IP patterns faster than Rust** | — |
| **Rust regex** | Aho-Corasick (any count), mature DFA, overall fastest | IP patterns slower than coregex |

**v0.9.5 Improvements:**
- 🚀 `multi_literal`: **99x faster** than stdlib (was 27x in v0.9.4!)
  - Literal extractor fix: factored prefixes now correctly expanded
  - `(Wanderlust|Weltanschauung)` now extracts full literals, not just `W`
  - Gap vs Rust reduced from 10x to 2.8x
- 🚀 `literal_alt`: **104x faster** than stdlib (was 81x)
- 🚀 `version`: **93x faster** than stdlib (was 70x)
- 🚀 Teddy Slim now supports 32 patterns (was 8)

**v0.9.4 Features:**
- ✅ `char_class`: **9x faster** than stdlib
- ✅ Teddy 2-byte fingerprint reduces false positives 90%

**v0.9.2 Features:**
- ✅ `ip`: **157x faster** than stdlib, **3.9x faster than Rust**
- ✅ Aho-Corasick for >8 literal patterns
- ✅ DigitPrefilter for digit-lead patterns

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
