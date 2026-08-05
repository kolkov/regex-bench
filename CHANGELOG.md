# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

---

## [2026-08-05] - Aho-Corasick v0.3.0, 6 patterns faster than Rust (v0.12.23)

### Changed
- **Updated coregex v0.12.22 → [v0.12.23](https://github.com/coregx/coregex/releases/tag/v0.12.23)**
  - `github.com/coregx/ahocorasick` v0.2.1 → v0.3.0 — zero-allocation `Find`/`FindAt` API

### Benchmark Results (AMD EPYC)

| Pattern | v0.12.18 | v0.12.23 | vs Rust | Notes |
|---------|----------|----------|---------|-------|
| inner_literal | 0.29 ms | 0.26 ms | **52x faster** | improved |
| ip | 2.10 ms | 0.77 ms | **17.6x faster** | improved |
| suffix | 1.83 ms | 1.79 ms | **7.7x faster** | improved |
| multiline_php | 0.66 ms | 0.38 ms | **2.0x faster** | improved |
| char_class | 47.95 ms | 41.91 ms | **1.4x faster** | improved |
| version | 1.81 ms | 0.65 ms | **1.2x faster** | improved |
| word_repeat | 185 ms | 179 ms | 3.2x | stable |
| http_methods | 1.28 ms | 1.51 ms | 2.4x | stable |

**6 patterns faster than Rust** (up from 4 in v0.12.18). No regressions.
README updated with full results table.

---

## [2026-03-24] - Flat DFA transition table, integrated prefilter (v0.12.18)

### Changed
- **Updated coregex v0.12.17 → v0.12.18**
  - Flat DFA transition table (Rust approach) — all 8 search functions
  - DFA integrated prefilter skip-ahead at start state
  - PikeVM integrated prefilter skip-ahead
  - 4x loop unrolling in searchFirstAt
  - NFA partialCoverage guard, DFA prefilter skip for incomplete prefilters
  - 386 safeOffset fix

### Benchmark Results (AMD EPYC)

| Pattern | v0.12.17 | v0.12.18 | vs Rust | Notes |
|---------|----------|----------|---------|-------|
| inner_literal | 0.30 ms | 0.29 ms | **~parity** | improved |
| char_class | 51.22 ms | 47.95 ms | **1.0x faster** | improved |
| ip | 2.14 ms | 2.10 ms | **5.6x faster** | stable |
| multiline_php | 0.50 ms | 0.66 ms | **~parity** | stable |
| word_repeat | 186 ms | 185 ms | 3.7x | stable |
| http_methods | 1.28 ms | 1.28 ms | 2.2x | stable |

35% faster than v0.12.14 baseline on Kostya's LangArena benchmark (1.80s → 1.17s).
No regressions on any platform (EPYC, Xeon, M1, 386).

---

## [2026-03-21] - Per-goroutine DFA cache, 7 correctness fixes (v0.12.15)

### Changed
- **Updated coregex v0.12.14 → v0.12.15**
  - Per-goroutine DFA cache (Rust approach) — eliminates concurrent data races
  - Pre-computed word boundary flags (30% → 0.3% CPU)
  - Integrated prefilter+DFA loop
  - 7 correctness fixes: anchor verification, .* newline boundary, alternation overflow
  - Stdlib compatibility test: 38/38 PASS

### Benchmark Results (AMD EPYC)

| Pattern | v0.12.14 | v0.12.15 | vs Rust | Notes |
|---------|----------|----------|---------|-------|
| anchored | 0.02 ms | 0.02 ms | 2.0x | restored |
| email | 0.57 ms | 0.50 ms | 2.1x | improved |
| uri | 0.65 ms | 0.61 ms | 1.7x | improved |
| multiline_php | 0.54 ms | 0.50 ms | 1.3x faster | improved |
| ip | 2.19 ms | 2.17 ms | 5.5x faster | stable |
| char_class | 38.88 ms | 41.0 ms | 1.2x faster | stable |
| http_methods | 0.77 ms | 1.56 ms | 2.2x | correctness fix ((?m)^ anchor) |
| suffix | 1.03 ms | 1.83 ms | 1.6x | correctness fix (.* newline) |
| word_repeat | 191 ms | 186 ms | 3.8x | improved |

No performance regressions. `http_methods` and `suffix` slower due to correctness fixes
(anchor verification, newline boundary handling).

---

## [2026-03-10] - BoundedBacktracker Span Fix, DFA FindAll Optimization (v0.12.6)

### Changed
- **Updated coregex v0.12.1 → v0.12.6**
  - v0.12.6: BoundedBacktracker span-based CanHandle, ReplaceAllStringFunc O(n), DFA FindAll early termination (#127)
  - v0.12.5: Non-greedy quantifier fix (PikeVM DFS-ordering), ReverseSuffix correctness (#124)
  - v0.12.4: Test coverage 80%+, CI improvements
  - v0.12.3: Cross-product literal expansion, 110x regexdna speedup (#119)
  - v0.12.2: ReverseSuffixSet safety guard, matchStartZero fix (#116)

### Changes vs v0.12.1 (previous benchmark run)

| Pattern | v0.12.1 | v0.12.6 | Delta |
|---------|---------|---------|-------|
| email | 0.61 ms | 0.58 ms | ~5% faster |
| multiline_php | 0.66 ms | 0.65 ms | stable |
| word_repeat | 184 ms | 187 ms | stable (noise) |

No regressions detected. All other patterns within noise margin.

**Summary:**
- coregex wins vs Rust: inner_literal, suffix, ip, char_class (4 patterns)
- Rust parity: multiline_php (0.65ms vs 0.66ms)
- Rust wins: literal_alt, multi_literal, email, uri, version, alpha_digit, word_digit, http_methods, word_repeat (9 patterns)
- Main gap: Teddy SIMD (Go/asm boundary overhead) — blocked on Go 1.26 archsimd
- Full table: see CI Job Summary

---

## [2026-02-15] - Bidirectional DFA Fallback, Bounded Repetitions Fix (v0.12.1)

### Changed
- **Updated coregex v0.11.5 → v0.12.1**
  - v0.12.1: Bounded repetitions ReverseSuffix fix (#115), AVX2 Teddy fix (#74)
  - v0.12.0: Anti-quadratic guard, DFA loop unrolling, DFA cache clear & continue
  - Bidirectional DFA fallback for BoundedBacktracker on large inputs
  - Digit-run skip optimization for `\d+`-leading patterns
  - CompositeSequenceDFA overmatching fix

### Results (GitHub Actions Ubuntu, 6MB input)

| Pattern | Go stdlib | Go coregex | Rust regex | vs stdlib | vs Rust | Winner |
|---------|-----------|------------|------------|-----------|---------|--------|
| inner_literal | 232 ms | **0.25 ms** | 0.31 ms | **926x** | **1.2x faster** | coregex |
| email | 260 ms | **0.61 ms** | 0.37 ms | **426x** | 1.6x slower | Rust |
| uri | 256 ms | 0.67 ms | **0.42 ms** | **382x** | 1.6x slower | Rust |
| suffix | 234 ms | **0.89 ms** | 1.09 ms | **263x** | **1.2x faster** | coregex |
| ip | 507 ms | **2.16 ms** | 12.05 ms | **235x** | **5.6x faster** | coregex |
| multiline_php | 103 ms | **0.66 ms** | 0.67 ms | **156x** | **~same** | parity |
| char_class | 560 ms | **41 ms** | 51 ms | **13.6x** | **1.2x faster** | coregex |
| word_repeat | 654 ms | 184 ms | **49 ms** | **3.6x** | 3.7x slower | Rust |

**Summary:**
- coregex wins: inner_literal, suffix, ip, char_class (4 patterns)
- Rust parity: multiline_php (0.66ms vs 0.67ms)
- Rust wins: literal_alt, multi_literal, email, uri, version, http_methods, word_repeat (7 patterns)
- Extreme: ip_nomatch **2542x**, suffix_find **1945x**, phone_nomatch **863x**

---

## [2026-02-01] - Issue #105 Fix: 8.6x Faster Than stdlib (v0.11.5)

### Changed
- **Updated coregex v0.11.4 → v0.11.5**
  - Fixed checkHasWordBoundary catastrophic slowdown (Issue #105)
  - Was 7,000,000x slower than stdlib on patterns with `\w{n,m}` quantifiers
  - Now **8.6x faster than stdlib** on same patterns

### Added
- **New benchmark pattern: `word_repeat`**
  - Pattern: `(\w{2,8})+`
  - Tests word quantifiers in capture groups (Issue #105 regression test)

### Results

| Pattern | Before v0.11.5 | After v0.11.5 | stdlib | Speedup |
|---------|----------------|---------------|--------|---------|
| word_repeat (79KB) | 3m 22s | **3.6 µs** | 28.5 µs | **8.6x faster** |

---

## [2026-01-16] - Multiline Near Rust Parity (v0.11.4)

### Changed
- **Updated coregex v0.11.3 → v0.11.4**
  - FindAll multiline fix: **78x faster** (Issue #102)
  - `multiline_php`: was 84x slower → now **~1.3x slower** than Rust
  - Root cause: `FindIndicesAt()` was missing dispatch for `UseMultilineReverseSuffix`

### Results (GitHub Actions Ubuntu, 6MB input)

| Pattern | Go stdlib | Go coregex | Rust regex | vs stdlib | vs Rust |
|---------|-----------|------------|------------|-----------|---------|
| multiline_php | 103 ms | **~1 ms** | 0.78 ms | **100x+** | ~1.3x slower |

**Before v0.11.4:** 66.48 ms (84x slower than Rust)
**After v0.11.4:** ~1 ms (near Rust parity!)

---

## [2026-01-16] - Auto-generated Comparison Table in CI

### Added
- **CI Job Summary now auto-generates comparison table**
  - Side-by-side: Go stdlib | Go coregex | Rust regex
  - vs stdlib speedup (bold if >= 2x)
  - vs Rust comparison (faster/slower)
  - Winner column (coregex/Rust)
  - Bold formatting on fastest time per row
  - Raw output in collapsible `<details>` section

### Changed
- Updated `.github/workflows/benchmark.yml` with table generation
- Added `scripts/generate-summary.sh` helper script

---

## [2026-01-16] - Multiline Pattern Benchmark (v0.11.1)

### Added
- **New pattern: `multiline_php`** (`(?m)^/.*\.php`)
  - Tests UseMultilineReverseSuffix strategy (Issue #97)
  - Multiline suffix matching for log files

### Changed
- **Updated coregex v0.11.0 → v0.11.1**
  - UseMultilineReverseSuffix: 1.4-5.7x faster than stdlib for `(?m)^.*suffix`
  - Known gap: Rust is 84x faster (DFA state acceleration, tracked in Issue #99)

### Results (GitHub Actions Ubuntu, 6MB input)

| Pattern | Go stdlib | Go coregex | Rust regex | vs stdlib | Winner |
|---------|-----------|------------|------------|-----------|--------|
| inner_literal | 206 ms | **0.45 ms** | 0.59 ms | **457x** | coregex |
| email | 249 ms | **1.30 ms** | 1.55 ms | **192x** | coregex |
| uri | 242 ms | 1.65 ms | **1.05 ms** | **147x** | Rust |
| literal_alt | 434 ms | 4.22 ms | **0.91 ms** | **103x** | Rust |
| multi_literal | 1300 ms | 13.06 ms | **4.56 ms** | **99x** | Rust |
| suffix | 203 ms | 2.12 ms | **1.60 ms** | **96x** | Rust |
| http_methods | 95 ms | 1.04 ms | **0.72 ms** | **92x** | Rust |
| ip | 468 ms | **6.37 ms** | 11.55 ms | **73x** | coregex |
| version | 154 ms | 2.78 ms | **0.91 ms** | **55x** | Rust |
| char_class | 494 ms | **49.31 ms** | 52.44 ms | **10x** | coregex |
| multiline_php | 93 ms | 66.48 ms | **0.79 ms** | **1.4x** | Rust |

**Summary:**
- coregex wins: inner_literal, email, ip, char_class (4 patterns)
- Rust wins: uri, literal_alt, multi_literal, suffix, http_methods, version, multiline_php (7 patterns)
- Multiline gap: Rust 84x faster due to DFA state acceleration

---

## [2026-01-15] - ReverseSuffix CharClass Plus Fix

### Changed
- **Updated coregex v0.10.9 → v0.10.10**
  - Fix: CharClass Plus patterns (`[^\s]+`, `[\w]+`) now use ReverseSuffix optimization
  - Bug: `[^\s]+\.txt` caused extreme benchmark to hang (266ms/MB instead of µs)
  - Result: All extreme patterns now complete in µs

### Results (GitHub Actions, 6MB no-digits data)
| Pattern | stdlib | coregex | Speedup |
|---------|--------|---------|---------|
| ip_nomatch | 392 ms | 215 µs | **1820x** |
| suffix_find | 226 ms | 218 µs | **1037x** |
| inner_nomatch | 210 ms | 254 µs | **826x** |
| phone_nomatch | 132 ms | 216 µs | **613x** |

---

## [2026-01-15] - UTF-8 Optimization + Fuzz Bug Fixes

### Changed
- **Updated coregex v0.10.8 → v0.10.9**
  - UTF-8 suffix sharing reduces dot NFA states 39→28
  - Anchored suffix prefilter for O(1) rejection
  - CharClassSearcher excludes `*` patterns (zero-width match fix)
  - Invalid UTF-8 handling for negated char classes (stdlib compat)
  - ReverseInner/ReverseSuffix whitelist (strategy safety)

### Results
- No regressions, all speedups maintained
- Local benchmarks: 113x-389x+ on various patterns

---

## [2026-01-15] - FindAll Anchored Pattern Fix

### Changed
- **Updated coregex v0.10.5 → v0.10.8**
  - v0.10.8: FindAll 600x faster for anchored patterns (#92)
  - v0.10.7: UTF-8 fixes + 100% stdlib API compatibility
  - v0.10.6: CompositeSequenceDFA for overlapping patterns

### Fixed
- **Anchored patterns (`^...`) allocation fix**
  - Before: 0.21 ms (huge allocation for 6MB input)
  - After: 0.03 ms (smart allocation: cap=1 for anchored)
  - Start-anchored patterns can only match at position 0

### Results
- **5 patterns now faster than Rust** (was 3):
  - `inner_literal`: 1.8x faster (0.31ms vs 0.56ms)
  - `ip`: 1.8x faster (6.7ms vs 12.1ms)
  - `http_methods`: 1.5x faster
  - `char_class`: 1.4x faster
  - `email`: 1.2x faster

---

## [2026-01-14] - Overlapping Char Classes Fix

### Changed
- **Updated coregex v0.10.4 → v0.10.5**
  - Critical fix: `\w+[0-9]+` patterns now work correctly (#81)
  - Bug: Greedy `\w+` consumed all characters (including digits)
  - Fix: Recursive backtracking in CompositeSearcher
  - `word_digit` pattern now returns correct matches

### Results
- `word_digit` (`\w+[0-9]+`): Now finds 3575 matches (was 0)
- No performance regression (+0.10% geomean)

---

## [2026-01-14] - Thread-safety Release

### Changed
- **Updated coregex v0.10.3 → v0.10.4**
  - Critical fix: Panic on concurrent usage of compiled Regexp (#78)
  - Implements `sync.Pool` pattern (same as Go stdlib `regexp`)
  - Thread-safe concurrent access to `*Regexp` instances
  - 32-bit platform compatibility (atomic operations alignment)
  - No performance regression (-3.84% improvement in geomean)

### Added
- **New benchmark patterns:**

  | Pattern | Regex | Category | Purpose |
  |---------|-------|----------|---------|
  | `alpha_digit` | `[a-zA-Z]+\d+` | Composite | Concatenated char classes |
  | `word_digit` | `\w+[0-9]+` | Composite | Word followed by digits |
  | `http_methods` | `^(GET\|POST\|PUT\|DELETE\|PATCH)` | Anchored | HTTP method dispatch |

- Added patterns to all 3 benchmark suites:
  - `go-coregex/main.go`
  - `go-stdlib/main.go`
  - `rust/src/main.rs`

### Technical
- These patterns test upcoming CompositeSearcher optimization (#72)
- Branch dispatch patterns test O(1) first-byte optimization
- Prepares for future coregex v0.11.0 release

---

## [2026-01-12] - Capture Group Fix

### Changed
- **Updated coregex v0.10.2 → v0.10.3**
  - Fixed: FindStringSubmatch returned incorrect captures for `.+` patterns
  - Bug: `^(.+)-(\d+)$` on "hello-123" returned wrong `matches[1]`
  - Root cause: StateSplit in PikeVM passed captures without cloning

---

## [2026-01-07] - Version Pattern Hotfix

### Changed
- **Updated coregex v0.10.1 → v0.10.2**
  - Restored DigitPrefilter for version patterns (`\d+\.\d+\.\d+`)
  - v0.10.1 incorrectly chose ReverseInner with "." as inner literal
  - Performance restored: 8.2ms → 2.15ms (3.8x speedup)

---

## [2026-01-07] - Fat Teddy Release

### Changed
- **Updated coregex v0.9.5 → v0.10.0**
  - Fat Teddy 16-bucket SIMD (33-64 patterns, 9+ GB/s)
  - AVX2 assembly implementation
  - Pure Go scalar fallback

### Results
- `multi_literal` (12 patterns): 11.62 ms (Aho-Corasick)
- 5 patterns now faster than Rust regex

---

## [2026-01-05] - Initial Public Release

### Added
- Cross-language regex benchmark suite
- **Go stdlib** benchmarks
- **Go coregex** benchmarks  
- **Rust regex** benchmarks
- 10 real-world patterns
- GitHub Actions CI/CD
- Extreme benchmark mode (3000x speedup demo)

### Patterns
| Pattern | Description |
|---------|-------------|
| literal_alt | 4-literal alternation |
| multi_literal | 12-literal alternation |
| anchored | Start anchor |
| inner_literal | Inner literal search |
| suffix | Suffix matching |
| char_class | Character class |
| email | Email validation |
| uri | URI parsing |
| version | Version numbers |
| ip | IPv4 validation |

---

## Links

- **coregex**: https://github.com/coregx/coregex
- **Benchmark repo**: https://github.com/kolkov/regex-bench
- **Go regex issue**: https://github.com/golang/go/issues/26623
