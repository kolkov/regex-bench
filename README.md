# regex-bench

Cross-language regex benchmark for **real-world patterns**.

Created to provide data for [golang/go#26623](https://github.com/golang/go/issues/26623) discussion on Go regex performance.

## Results

**Intel i7-1255U, 6.0 MB input text**

| Pattern | Go stdlib | Go coregex | Rust regex | Best Go vs stdlib |
|---------|-----------|------------|------------|-------------------|
| literal_alt | 582 ms | 42 ms | **11 ms** | **14x faster** |
| anchored | <1 ms | <1 ms | <1 ms | — |
| inner_literal | 281 ms | 3 ms | **1.3 ms** | **94x faster** |
| suffix | 289 ms | **2 ms** | 1.9 ms | **145x faster** |
| char_class | 617 ms | **28 ms** | 66 ms | **22x faster** |
| email | 318 ms | **1.9 ms** | 3.4 ms | **167x faster** |

### Key Findings

**Go coregex v0.8.22 vs Go stdlib:**
- `email`: **167x faster**
- `suffix`: **145x faster**  
- `inner_literal`: **94x faster**
- `char_class`: **22x faster**
- `literal_alt`: **14x faster**

**Go coregex vs Rust regex:**
- `char_class`: **coregex 2.4x faster** (28ms vs 66ms)
- `email`: **coregex 1.8x faster** (1.9ms vs 3.4ms)
- `suffix`: **comparable** (2ms vs 1.9ms)
- `inner_literal`: Rust 2.3x faster
- `literal_alt`: Rust 4x faster (Aho-Corasick)

### Analysis

Go's stdlib `regexp` uses a simple NFA without prefilters or optimizations. Both coregex and Rust's regex use:
- Lazy DFA with on-demand state compilation
- SIMD prefilters (Teddy, memchr)
- Reverse search strategies for `.*` patterns
- Specialized searchers (CharClassSearcher in coregex)

Rust's advantage on `literal_alt` comes from Aho-Corasick integration. coregex wins on character classes due to CharClassSearcher's 256-byte lookup table.

## Patterns Tested

| Pattern | Regex | Type |
|---------|-------|------|
| literal_alt | `error\|warning\|fatal\|critical` | Multi-literal |
| anchored | `^HTTP/[12]\.[01]` | Start anchor |
| inner_literal | `.*@example\.com` | Inner literal |
| suffix | `.*\.(txt\|log\|md)` | Suffix match |
| char_class | `[\w]+` | Character class |
| email | `[\w.+-]+@[\w.-]+\.[\w.-]+` | Complex |

## Running Benchmarks

```bash
# Generate input data (6 MB)
go run scripts/generate-input.go

# Build and run Go
cd go-stdlib && go build -o ../bin/go-stdlib.exe . && cd ..
cd go-coregex && go build -o ../bin/go-coregex.exe . && cd ..
./bin/go-stdlib.exe input/data.txt
./bin/go-coregex.exe input/data.txt

# Run Rust (requires WSL or Linux)
wsl ./bin/rust-benchmark input/data.txt
```

## Links

- **coregex**: https://github.com/coregx/coregex
- **Go issue**: https://github.com/golang/go/issues/26623
- **Rust regex**: https://github.com/rust-lang/regex

## License

MIT
