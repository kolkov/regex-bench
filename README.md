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

**Intel i7-1255U, 6.0 MB input text**

| Pattern | Go stdlib | Go coregex | Rust regex | Winner |
|---------|-----------|------------|------------|--------|
| literal_alt | 369 ms | 32 ms | **4.3 ms** | Rust **86x** |
| anchored | <1 ms | <1 ms | 0.6 ms | — |
| inner_literal | 172 ms | 604 ms | **0.9 ms** | Rust **191x** |
| suffix | 183 ms | 393 ms | **1.5 ms** | Rust **122x** |
| char_class | 460 ms | 500 ms | **62 ms** | Rust **7x** |
| email | 212 ms | 396 ms | **1.8 ms** | Rust **118x** |

### Key Insights

- **Rust regex** is the gold standard - 7x to 200x faster than Go on all patterns
- **Go coregex** excels at literal alternations (11.5x faster than stdlib)
- **Go stdlib** is better on `.*` patterns than coregex (optimized NFA)
- Huge optimization potential for Go engines on `.*` patterns

### Go-only Comparison

| Pattern | Go stdlib | Go coregex | Winner |
|---------|-----------|------------|--------|
| literal_alt | 369 ms | **32 ms** | coregex **11.5x** |
| inner_literal | **172 ms** | 604 ms | stdlib 3.5x |
| suffix | **183 ms** | 393 ms | stdlib 2.1x |
| char_class | **460 ms** | 500 ms | stdlib 1.1x |
| email | **212 ms** | 396 ms | stdlib 1.9x |

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
