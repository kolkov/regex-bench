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
| literal_alt | 482 ms | 4.3 ms | **0.7 ms** | **112x** | 5.9x slower |
| multi_literal | 1402 ms | 12.4 ms | **4.8 ms** | **113x** | 2.6x slower |
| anchored | 0.05 ms | **0.02 ms** | 0.04 ms | **2.5x** | **2x faster** |
| inner_literal | 232 ms | **0.83 ms** | 0.54 ms | **280x** | 1.5x slower |
| suffix | 233 ms | **1.0 ms** | 1.3 ms | **224x** | **27% faster** |
| char_class | 511 ms | **39.4 ms** | 53.4 ms | **13x** | **35% faster** |
| email | 261 ms | **1.2 ms** | 1.4 ms | **225x** | **16% faster** |
| uri | 258 ms | 1.8 ms | **1.2 ms** | **145x** | 1.5x slower |
| version | 171 ms | 2.2 ms | **0.7 ms** | **79x** | 3.2x slower |
| **ip** | 500 ms | **3.8 ms** | 12.3 ms | **132x** | **3.3x faster** |

### Key Findings

**Go coregex v0.10.0 vs Go stdlib:**
- Most patterns: **13-280x faster**
- Best: `inner_literal` **280x**, `email` **225x**, `suffix` **224x**, `ip` **132x**
- `multi_literal` **113x** (Aho-Corasick)
- `literal_alt` **112x** (Teddy SIMD)
- `char_class` **13x** (CharClassSearcher)
- `version` **79x** (DigitPrefilter)

**Go coregex faster than Rust (5 patterns!):**
- `char_class`: **coregex 35% faster** (39ms vs 53ms)
- `ip`: **coregex 3.3x faster** (3.8ms vs 12.3ms)
- `suffix`: **coregex 27% faster** (1.0ms vs 1.3ms)
- `email`: **coregex 16% faster** (1.2ms vs 1.4ms)
- `anchored`: **coregex 2x faster** (0.02ms vs 0.04ms)

**Rust faster than coregex:**
- `literal_alt`: Rust 5.9x faster
- `version`: Rust 3.2x faster
- `multi_literal`: Rust 2.6x faster
- `inner_literal`: Rust 1.5x faster
- `uri`: Rust 1.5x faster

> **Note**: Rust regex has 10+ years of development. coregex optimizations are targeted, not universal.

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 13-225x slower |
| **Go coregex** | Reverse search, SIMD prefilters, Aho-Corasick, **5 patterns faster than Rust** | — |
| **Rust regex** | Aho-Corasick (any count), mature DFA, overall fastest | char_class, IP, suffix, email slower |

**v0.10.0 Improvements:**
- Fat Teddy AVX2: 33-64 pattern support (9+ GB/s throughput)
- **5 patterns now faster than Rust**: char_class, ip, suffix, email, anchored
- `inner_literal`: **280x faster** (was 140x), gap vs Rust reduced from 2.3x to 1.5x
- `suffix`: **224x faster** (was 158x), now **27% faster than Rust**
- `char_class`: **13x faster** (was 9x), now **35% faster than Rust**
- `email`: **225x faster** (was 145x), now **16% faster than Rust**

**v0.9.5 Improvements:**
- `multi_literal`: **99x faster** than stdlib (was 27x in v0.9.4!)
  - Literal extractor fix: factored prefixes now correctly expanded
  - `(Wanderlust|Weltanschauung)` now extracts full literals, not just `W`
  - Gap vs Rust reduced from 10x to 2.8x
- `literal_alt`: **104x faster** than stdlib (was 81x)
- Teddy Slim now supports 32 patterns (was 8)

**v0.9.4 Features:**
- `char_class`: **9x faster** than stdlib
- Teddy 2-byte fingerprint reduces false positives 90%

**v0.9.2 Features:**
- `ip`: **157x faster** than stdlib, **3.9x faster than Rust**
- Aho-Corasick for >8 literal patterns
- DigitPrefilter for digit-lead patterns

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
