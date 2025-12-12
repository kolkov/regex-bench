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

**Intel i7-1255U, 6.0 MB input text** | coregex v0.8.19

| Pattern | Go stdlib | Go coregex | Rust regex | Winner |
|---------|-----------|------------|------------|--------|
| literal_alt | 810 ms | 50 ms | **13 ms** | Rust **62x** |
| anchored | 1 ms | <1 ms | 0.5 ms | — |
| inner_literal | 421 ms | 3 ms | **1.2 ms** | Rust **351x** |
| suffix | 393 ms | 826 ms | **1.9 ms** | Rust **206x** |
| char_class | 772 ms | 876 ms | **57 ms** | Rust **13x** |
| email | 395 ms | 831 ms | **1.9 ms** | Rust **207x** |

### Key Insights

- **Rust regex** is the gold standard - 13x to 350x faster than Go stdlib
- **Go coregex** excels at literal patterns:
  - `inner_literal` (`.*@example\.com`): **140x faster** than stdlib
  - `literal_alt` (`error|warning|...`): **16x faster** than stdlib
- **Go stdlib** is better on suffix alternation and email patterns
- Suffix alternation (`.*\.(txt|log)`) needs optimization in coregex

### Go-only Comparison

| Pattern | Go stdlib | Go coregex | Winner |
|---------|-----------|------------|--------|
| literal_alt | 810 ms | **50 ms** | coregex **16x** |
| inner_literal | 421 ms | **3 ms** | coregex **140x** |
| suffix | **393 ms** | 826 ms | stdlib 2.1x |
| char_class | **772 ms** | 876 ms | stdlib 1.1x |
| email | **395 ms** | 831 ms | stdlib 2.1x |

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
