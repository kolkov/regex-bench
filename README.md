# regex-bench

Cross-language regex benchmark for **real-world patterns**.

Created to provide data for [golang/go#26623](https://github.com/golang/go/issues/26623) discussion on Go regex performance.

## Test Environment

All benchmarks run on **identical conditions**:
- **OS**: Linux (Ubuntu via WSL2 or GitHub Actions)
- **Input**: 6.0 MB generated text file
- **Method**: Each engine compiled natively, same input file, same patterns

> **Note**: Cross-compiled Go binaries run in WSL2 for fair comparison with Rust.

## Results ([v0.12.23](https://github.com/coregx/coregex/releases/tag/v0.12.23))

**GitHub Actions Ubuntu (AMD EPYC), 6.0 MB input** (using `FindAll` for fair comparison)

| Pattern | Go stdlib | Go coregex | Rust regex | vs stdlib | vs Rust | Winner |
|---------|-----------|------------|------------|-----------|---------|--------|
| inner_literal | 232 ms | **0.26 ms** | 13.55 ms | **893x** | **52x faster** | coregex |
| ip | 489 ms | **0.77 ms** | 13.53 ms | **635x** | **17.6x faster** | coregex |
| email | 257 ms | **0.55 ms** | 0.26 ms | **467x** | 2.1x slower | Rust |
| uri | 258 ms | **0.61 ms** | 0.34 ms | **424x** | 1.8x slower | Rust |
| multiline_php | 101 ms | **0.38 ms** | 0.76 ms | **266x** | **2.0x faster** | coregex |
| version | 163 ms | **0.65 ms** | 0.79 ms | **250x** | **1.2x faster** | coregex |
| suffix | 236 ms | **1.79 ms** | 13.70 ms | **132x** | **7.7x faster** | coregex |
| literal_alt | 232 ms | **4.69 ms** | 0.63 ms | **49x** | 7.4x slower | Rust |
| http_methods | 103 ms | **1.51 ms** | 0.64 ms | **68x** | 2.4x slower | Rust |
| multi_literal | 236 ms | **12.89 ms** | 5.32 ms | **18x** | 2.4x slower | Rust |
| char_class | 507 ms | **41.91 ms** | 58.38 ms | **12x** | **1.4x faster** | coregex |
| alpha_digit | 255 ms | 29.02 ms | **13.48 ms** | **9x** | 2.2x slower | Rust |
| word_digit | 269 ms | 29.22 ms | **13.59 ms** | **9x** | 2.2x slower | Rust |
| word_repeat | 647 ms | 179 ms | **56 ms** | **3.6x** | 3.2x slower | Rust |
| anchored | 0.04 ms | 0.05 ms | 0.04 ms | ~1x | ~same | — |
| anchored_php | 0.05 ms | 0.06 ms | 0.38 ms | ~1x | ~same | — |

> **coregex [v0.12.23](https://github.com/coregx/coregex/releases/tag/v0.12.23)** — Multi-engine architecture with 17 strategies, SIMD prefilters, zero-allocation Aho-Corasick v0.3.0. Run `make extreme` for 2500x demo.

### Key Findings

**Go coregex v0.12.23 vs Go stdlib:**
- All patterns: **3.6-893x faster**
- Best: `inner_literal` **893x**, `ip` **635x**, `email` **467x**
- `uri` **424x**, `multiline_php` **266x**, `version` **250x**
- `suffix` **132x**, `http_methods` **68x**, `literal_alt` **49x**
- `char_class` **12x** (CharClassSearcher, **faster than Rust!**)
- `word_repeat` **3.6x** (flat DFA with 4x unrolling)

**Go coregex faster than Rust (6 patterns):**
- `inner_literal`: **coregex 52x faster** (0.26ms vs 13.55ms)
- `ip`: **coregex 17.6x faster** (0.77ms vs 13.53ms)
- `suffix`: **coregex 7.7x faster** (1.79ms vs 13.70ms)
- `multiline_php`: **coregex 2.0x faster** (0.38ms vs 0.76ms)
- `char_class`: **coregex 1.4x faster** (41.9ms vs 58.4ms)
- `version`: **coregex 1.2x faster** (0.65ms vs 0.79ms)

**Rust faster than coregex:**
- `literal_alt`: Rust 7.4x faster (Teddy with more buckets)
- `word_repeat`: Rust 3.2x faster (DFA state acceleration)
- `multi_literal`: Rust 2.4x faster
- `http_methods`: Rust 2.4x faster
- `alpha_digit`, `word_digit`: Rust 2.2x faster
- `email`: Rust 2.1x faster
- `uri`: Rust 1.8x faster

> **Note**: Rust regex has 10+ years of development. coregex optimizations are targeted, not universal.

### Analysis

| Engine | Strengths | Weaknesses |
|--------|-----------|------------|
| **Go stdlib** | Simple, no dependencies | No optimizations, 3.6-893x slower |
| **Go coregex** | Flat DFA, reverse search, SIMD prefilters, Aho-Corasick, bidirectional DFA, **6 patterns faster than Rust** | Teddy Go/ASM gap, word_repeat |
| **Rust regex** | DFA state acceleration, Teddy Fat, mature DFA | ip, inner_literal, suffix, multiline_php, char_class, version slower than coregex |

**v0.12.23 (Current):**
- 17 strategies: NFA, DFA, OnePass, BoundedBacktracker, Teddy, Aho-Corasick, reverse search, and more
- Flat DFA transition table (Rust approach) — single flat array, no pointer chase
- SIMD prefilters: AVX2 memchr, SSSE3/AVX2 Teddy, Aho-Corasick DFA
- Zero-allocation API: `Find`/`FindAt`/`IsMatch` with no heap allocations
- **6 patterns faster than Rust**: inner_literal (52x), ip (17.6x), suffix (7.7x), multiline_php (2.0x), char_class (1.4x), version (1.2x)

**Historical Improvements:**
- v0.12.23: Aho-Corasick v0.3.0 (zero-alloc Find/FindAt API)
- v0.12.22: Lazy memory architecture — 5-7x memory reduction per pattern
- v0.12.18: Flat DFA transition table, integrated prefilter, 4x unrolling — 3x from Rust
- v0.12.17: Fix LogParser ARM64 regression, restore DFA/Teddy for (?m)^
- v0.12.16: WrapLineAnchor for (?m)^ patterns
- v0.12.15: Per-goroutine DFA cache, 7 correctness fixes, stdlib compat test (38/38)
- v0.12.14: Concurrent isMatchDFA safety fix (#137)
- v0.12.13: FatTeddy AVX2 fix, prefilter acceleration, AC v0.2.1
- v0.12.1: Bidirectional DFA fallback, bounded repetitions fix (#115), AVX2 Teddy fix (#74)
- v0.12.0: Anti-quadratic guard, DFA loop unrolling, DFA cache clear & continue
- v0.11.4: FindAll multiline fix, 78x faster (Issue #102)
- v0.11.3: UseMultilineReverseSuffix prefix fast path 319-552x (Issue #99)
- v0.11.1: UseMultilineReverseSuffix for multiline patterns (Issue #97)
- v0.11.0: UseAnchoredLiteral 32-133x speedup (Issue #79)
- v0.10.10: ReverseSuffix CharClass Plus fix
- v0.10.9: UTF-8 optimization + fuzz-found bug fixes
- v0.10.8: FindAll allocation fix for anchored patterns
- v0.10.7: UTF-8 fixes + 100% stdlib API compatibility
- v0.10.5: CompositeSearcher backtracking fix
- v0.10.0: Fat Teddy AVX2 (33-64 patterns, 9+ GB/s)
- v0.9.5: Aho-Corasick integration, Teddy 32 patterns

## Extreme Speedups (1000-3000x)

The "3-3000x faster" claim refers to **specific edge cases** where coregex prefilters can skip entire input:

```bash
make extreme       # Run on no-match data (~300-560x)
make extreme-3000x # Run on no-digits data (1000-3000x)
```

**GitHub Actions Ubuntu results** (6 MB no-digits data, v0.12.1):

| Pattern | Go stdlib | Go coregex | Speedup |
|---------|-----------|------------|---------|
| ip_nomatch | 422 ms | 166 µs | **2542x** |
| suffix_find | 245 ms | 126 µs | **1945x** |
| phone_nomatch | 143 ms | 166 µs | **863x** |
| inner_nomatch | 229 ms | 382 µs | **598x** |

[![Extreme Benchmark](https://github.com/kolkov/regex-bench/actions/workflows/extreme-benchmark.yml/badge.svg)](https://github.com/kolkov/regex-bench/actions/workflows/extreme-benchmark.yml)

> **Note**: Results vary between runs (±30%) due to CI VM load and OS scheduling.
> The key insight: coregex operates in **microseconds**, stdlib in **hundreds of milliseconds**.

**When do we see 3000x?**

The 3000x speedup occurs in coregex's own benchmark suite (`go test -bench`) under specific conditions:
- **Pattern**: IP regex on data with NO IP addresses
- **Size**: 1 MB of pure text
- **Measurement**: `go test -bench` with multiple iterations

```go
// In coregex repo:
BenchmarkIPRegex_Find/stdlib_1MB_no_ips    74.5ms
BenchmarkIPRegex_Find/coregex_1MB_no_ips   22.4µs  // 3324x
```

The extreme speedup happens because:
1. **DigitPrefilter** scans for first digit character
2. No digits in input → entire 1 MB skipped in ~20µs
3. stdlib must scan byte-by-byte → 74ms

**Verified speedups** (from coregex repo, `docs/dev/SPEEDUP_VERIFICATION.md`):

| Pattern | Strategy | Max Speedup |
|---------|----------|-------------|
| IP no-match (1MB) | DigitPrefilter | **3324x** |
| `.*\.txt$` (1MB) | ReverseSuffix | **1124x** |
| `.*error.*` (32KB) | ReverseInner | **909x** |

> The speedup depends on input characteristics. Real-world mixed data shows 15-560x.

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
| http_methods | `(?m)^(GET\|POST\|PUT\|DELETE\|PATCH)` | Multiline log parsing | BranchDispatch |
| anchored_php | `^/.*[\w-]+\.php` | URL path matching | UseAnchoredLiteral |
| multiline_php | `(?m)^/.*\.php` | Multiline PHP paths | UseMultilineReverseSuffix |
| word_repeat | `(\w{2,8})+` | Word quantifiers | BoundedBacktracker + DFA fallback |

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

**Auto-generated comparison table** in Job Summary:
- Side-by-side results for all 3 engines
- Speedup calculations (vs stdlib, vs Rust)
- Winner column with bold formatting
- Raw output in collapsible section

## Links

- **coregex**: https://github.com/coregx/coregex
- **Go issue**: https://github.com/golang/go/issues/26623
- **Rust regex**: https://github.com/rust-lang/regex

## License

MIT
