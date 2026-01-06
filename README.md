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
| literal_alt | 472 ms | 31 ms | **0.8 ms** | **15x faster** |
| multi_literal | 1405 ms | 43 ms | **4.8 ms** | **33x faster** |
| anchored | 0.03 ms | 0.09 ms | 0.03 ms | — |
| inner_literal | 231 ms | **1.4 ms** | 0.5 ms | **165x faster** |
| suffix | 234 ms | **1.4 ms** | 1.3 ms | **167x faster** |
| char_class | 514 ms | 139 ms | **53 ms** | **3.7x faster** |
| email | 259 ms | 2.5 ms | **1.6 ms** | **104x faster** |
| uri | 257 ms | 2.1 ms | **0.9 ms** | **122x faster** |
| version | 169 ms | 8.2 ms | **0.7 ms** | **21x faster** |
| **ip** | 496 ms | **3.9 ms** | 12.3 ms | **127x faster** |

### Key Findings

**Go coregex v0.9.2 vs Go stdlib:**
- Most patterns: **15-167x faster**
- Best: `suffix` **167x**, `inner_literal` **165x**, `ip` **127x**, `uri` **122x**
- `multi_literal` **33x** (Aho-Corasick for 12 patterns)
- `version` **21x** (ReverseInner)

**Go coregex vs Rust regex:**
- `ip`: coregex faster (3.9ms vs 12.3ms) — specific to this pattern
- `suffix`: ~tie (1.4ms vs 1.3ms)
- `inner_literal`: Rust 2.8x faster
- `email`: Rust 1.6x faster
- `multi_literal`: Rust 9x faster (optimized Aho-Corasick)
- `version`: Rust 12x faster
- `char_class`: Rust 2.6x faster
- `literal_alt`: Rust 39x faster (Aho-Corasick for any count)

> **Note**: Rust regex has 10+ years of development. coregex optimizations are targeted, not universal.

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 4-167x slower |
| **Go coregex** | Reverse search, SIMD prefilters, Aho-Corasick, **IP patterns faster than Rust** | Complex alternations |
| **Rust regex** | Aho-Corasick (any count), mature DFA, overall fastest | IP patterns slower than coregex |

**v0.9.2 Improvements:**
- 🚀 `ip`: **127x faster** than stdlib, **3.1x faster than Rust**
  - Compile-time strategy selection based on NFA complexity
  - Complex digit patterns (>100 NFA states) use optimized lazy DFA
  - Removed runtime overhead from v0.9.1 adaptive switching

**v0.9.0-v0.9.1 Features:**
- ✅ `multi_literal`: Aho-Corasick for >8 literal patterns
- ✅ `version`: ReverseInner with `.` literal
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
