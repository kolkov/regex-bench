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
| literal_alt | 374 ms | 35 ms | **0.8 ms** | **11x faster** |
| **multi_literal** | 1152 ms | **48 ms** | 4.9 ms | **24x faster** |
| anchored | 0.03 ms | 0.02 ms | **0.05 ms** | — |
| inner_literal | 224 ms | **1.9 ms** | 0.6 ms | **118x faster** |
| suffix | 203 ms | **2.0 ms** | 1.3 ms | **102x faster** |
| char_class | 502 ms | 147 ms | **53 ms** | **3.4x faster** |
| email | 244 ms | 10 ms | **1.5 ms** | **24x faster** |
| uri | 231 ms | 3.1 ms | **1.0 ms** | **75x faster** |
| **version** | 130 ms | **10 ms** | 0.65 ms | **13x faster** |
| **ip** | 478 ms | **137 ms** | **12 ms** | **3.5x faster** 🔧 |

### Key Findings

**Go coregex v0.9.1 vs Go stdlib:**
- Most patterns: **3-118x faster**
- Best: `inner_literal` **118x**, `suffix` **102x**, `uri` **75x**
- `multi_literal` **24x** (Aho-Corasick for 12 patterns)
- `version` **13x** (ReverseInner with `.` literal)
- 🔧 `ip` **3.5x** (fixed: was 1.3x regression in v0.9.0)

**Go coregex vs Rust regex:**
- `suffix`: **coregex ~tie** (1.2ms vs 1.3ms)
- `inner_literal`: Rust 2x faster
- `email`: Rust ~tie (1.9ms vs 1.5ms)
- `multi_literal`: Rust 9x faster (Aho-Corasick)
- `version`: Rust 13x faster
- `char_class`: Rust 2.6x faster
- `literal_alt`: Rust 39x faster (Aho-Corasick for any count)

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 3-199x slower |
| **Go coregex** | Reverse search, SIMD prefilters, Aho-Corasick (>8 patterns), DigitPrefilter | Complex nested alternations (ip) |
| **Rust regex** | Aho-Corasick (any count), mature DFA, overall fastest | — |

**v0.9.1 Fixes:**
- 🔧 `ip`: Fixed DigitPrefilter regression (was 1.3x slower, now 3.5x faster)
  - Runtime adaptive switching: after 64 consecutive false positives, switch to lazy DFA
  - Based on Rust regex insight: "prefilter with high FP rate makes search slower"

**v0.9.0 Features:**
- ✅ `multi_literal`: Aho-Corasick triggers for >8 literal patterns
- ✅ `version`: ReverseInner with `.` literal

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
| ip | `(?:(?:25[0-5]\|2[0-4][0-9]\|...)\.){3}...` | IPv4 validation | **UseBoth** (lazy DFA) |

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
