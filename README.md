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
| inner_literal | 231 ms | **0.68 ms** | 0.51 ms | **340x** | 33% slower |
| email | 277 ms | **1.28 ms** | 1.40 ms | **216x** | **9% faster** |
| suffix | 234 ms | 1.99 ms | **1.30 ms** | **118x** | 53% slower |
| uri | 267 ms | 1.85 ms | **0.91 ms** | **144x** | 2x slower |
| ip | 507 ms | **3.88 ms** | 12.30 ms | **131x** | **3.2x faster** |
| multi_literal | 1413 ms | 12.56 ms | **4.87 ms** | **112x** | 2.6x slower |
| literal_alt | 480 ms | 4.33 ms | **0.75 ms** | **111x** | 5.8x slower |
| char_class | 516 ms | **41.40 ms** | 53.76 ms | **12x** | **23% faster** |
| version | 177 ms | 8.21 ms | **0.71 ms** | **22x** | 11.6x slower |
| anchored | 0.03 ms | 0.21 ms | **0.05 ms** | 7x slower | 4x slower |

> **Note**: v0.10.1 has regressions on `version` (3.8x) and `anchored` (10x) patterns compared to v0.10.0.
> See [coregex issue #74](https://github.com/coregx/coregex/issues/74) for investigation.

### Key Findings

**Go coregex v0.10.1 vs Go stdlib:**
- Most patterns: **12-340x faster**
- Best: `inner_literal` **340x**, `email` **216x**, `uri` **144x**, `ip` **131x**
- `multi_literal` **112x** (Aho-Corasick)
- `literal_alt` **111x** (Teddy SIMD)
- `char_class` **12x** (CharClassSearcher)
- `version` **22x** (ReverseInner) — regression from v0.10.0 (was 79x)
- `anchored` **7x slower** — regression from v0.10.0 (was 2.5x faster)

**Go coregex faster than Rust (3 patterns):**
- `ip`: **coregex 3.2x faster** (3.9ms vs 12.3ms)
- `char_class`: **coregex 23% faster** (41ms vs 54ms)
- `email`: **coregex 9% faster** (1.3ms vs 1.4ms)

**Rust faster than coregex:**
- `version`: Rust 11.6x faster (regression!)
- `literal_alt`: Rust 5.8x faster
- `anchored`: Rust 4x faster (regression!)
- `multi_literal`: Rust 2.6x faster
- `uri`: Rust 2x faster
- `suffix`: Rust 53% faster
- `inner_literal`: Rust 33% faster

> **Note**: Rust regex has 10+ years of development. coregex optimizations are targeted, not universal.

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 13-225x slower |
| **Go coregex** | Reverse search, SIMD prefilters, Aho-Corasick, **3 patterns faster than Rust** | v0.10.1 regressions on version/anchored |
| **Rust regex** | Aho-Corasick (any count), mature DFA, overall fastest | char_class, IP, suffix, email slower |

**v0.10.1 Regressions (under investigation):**
- `version`: 2.2ms → 8.2ms (3.8x regression) — strategy changed from DigitPrefilter to ReverseInner
- `anchored`: 0.02ms → 0.21ms (10x regression) — same UseNFA strategy, cause unknown
- See [coregex #74](https://github.com/coregx/coregex/issues/74) for details

**v0.10.0 Improvements:**
- Fat Teddy AVX2: 33-64 pattern support (9+ GB/s throughput)
- **5 patterns faster than Rust** (before v0.10.1 regressions): char_class, ip, suffix, email, anchored
- `inner_literal`: **280x faster** (was 140x), gap vs Rust reduced from 2.3x to 1.5x
- `suffix`: **224x faster** (was 158x), was **27% faster than Rust**
- `char_class`: **13x faster** (was 9x), still **35% faster than Rust**
- `email`: **225x faster** (was 145x), still **16% faster than Rust**

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
| version | `\d+\.\d+\.\d+` | Version numbers | ReverseInner (regression from DigitPrefilter) |
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
