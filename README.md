# regex-bench

Cross-language regex benchmark focusing on **real-world patterns** where different engines excel.

## Philosophy

Unlike other benchmarks that use a single set of patterns, we test **multiple categories** to show where each engine shines:

| Category | Example | Tests |
|----------|---------|-------|
| **Literals** | `error\|warning\|fatal` | Literal extraction, multi-pattern |
| **Anchored** | `^HTTP/\d\.\d` | Start anchor optimization |
| **Inner Literal** | `.*@example\.com` | Reverse/inner literal search |
| **Suffix** | `.*\.(txt\|log)` | Suffix-based search |
| **Character Class** | `[\w]+` | Pure NFA/DFA performance |
| **Complex** | Email pattern | Real-world combined patterns |

## Results

**Intel i7-1255U, 6.0 MB input text** | coregex v0.8.20

| Pattern | Go stdlib | Go coregex | Rust regex | Winner |
|---------|-----------|------------|------------|--------|
| literal_alt | 760 ms | 43 ms | **13 ms** | Rust **58x** |
| anchored | <1 ms | <1 ms | 0.5 ms | — |
| inner_literal | 378 ms | 3 ms | **1.2 ms** | Rust **315x** |
| suffix | 372 ms | **5 ms** | 1.9 ms | coregex **74x** |
| char_class | 696 ms | 851 ms | **57 ms** | Rust **12x** |
| email | 368 ms | 758 ms | **1.9 ms** | Rust **193x** |

### Key Insights

- **Rust regex** is the gold standard - 12x to 315x faster than Go stdlib
- **Go coregex v0.8.20** excels at literal and suffix patterns:
  - `inner_literal` (`.*@example\.com`): **125x faster** than stdlib
  - `suffix` (`.*\.(txt|log|md)`): **74x faster** than stdlib - **NEW in v0.8.20!**
  - `literal_alt` (`error|warning|...`): **17.6x faster** than stdlib
- **ReverseSuffixSet optimization** (v0.8.20) - novel optimization NOT present in rust-regex!
- Character class and email patterns still need optimization

### Go-only Comparison

| Pattern | Go stdlib | Go coregex | Winner |
|---------|-----------|------------|--------|
| literal_alt | 760 ms | **43 ms** | coregex **17.6x** |
| inner_literal | 378 ms | **3 ms** | coregex **125x** |
| suffix | 372 ms | **5 ms** | coregex **74x** |
| char_class | **696 ms** | 851 ms | stdlib 1.2x |
| email | **368 ms** | 758 ms | stdlib 2.1x |

## Patterns Tested

```
literal_alt     error|warning|fatal|critical
anchored        ^HTTP/[12]\.[01]
inner_literal   .*@example\.com
suffix          .*\.(txt|log|md)
char_class      [\w]+
email           [\w.+-]+@[\w.-]+\.[\w.-]+
```

## Running Benchmarks

```bash
# Generate input data (6 MB)
go run scripts/generate-input.go

# Build Go
cd go-stdlib && go build -o ../bin/go-stdlib.exe .
cd go-coregex && go build -o ../bin/go-coregex.exe .

# Build Rust (via Docker)
docker run --rm -v "$(pwd):/app" -w /app/rust rust:latest \
  bash -c "cargo build --release && cp target/release/benchmark /app/bin/rust-benchmark"

# Run
./bin/go-stdlib.exe input/data.txt
./bin/go-coregex.exe input/data.txt
docker run --rm -v "$(pwd):/app" -w /app rust:latest ./bin/rust-benchmark input/data.txt
```

## Project Structure

```
regex-bench/
├── go-stdlib/      # Go standard library regexp
├── go-coregex/     # coregex high-performance engine
├── rust/           # Rust regex crate
├── input/          # Generated test data
├── scripts/        # Helper scripts
└── bin/            # Compiled binaries
```

## Adding a New Implementation

1. Create folder: `{language}-{library}/`
2. Implement benchmark following existing examples
3. Output format: `{pattern_name} {time_ms} ms {count} matches`
4. Submit PR

## License

MIT
