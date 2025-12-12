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

**Intel i7-1255U, Windows 11, 6.0 MB input text**

| Pattern | Go stdlib | Go coregex | Winner |
|---------|-----------|------------|--------|
| literal_alt | 369 ms | **32 ms** | coregex **11.5x** |
| anchored | <1 ms | <1 ms | — |
| inner_literal | **172 ms** | 604 ms | stdlib 3.5x |
| suffix | **183 ms** | 393 ms | stdlib 2.1x |
| char_class | **460 ms** | 500 ms | stdlib 1.1x |
| email | **212 ms** | 396 ms | stdlib 1.9x |

### Key Insights

- **coregex excels** at literal alternations (UseTeddy multi-pattern SIMD)
- **stdlib excels** at patterns starting with `.*` (optimized NFA)
- Character class performance is comparable

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

# Build
cd go-stdlib && go build -o ../bin/go-stdlib.exe .
cd go-coregex && go build -o ../bin/go-coregex.exe .

# Run
./bin/go-stdlib.exe input/data.txt
./bin/go-coregex.exe input/data.txt
```

Or use Make:

```bash
make all
```

## Project Structure

```
regex-bench/
├── go-stdlib/      # Go standard library regexp
├── go-coregex/     # coregex high-performance engine
├── rust/           # Rust regex (planned)
├── input/          # Generated test data
├── scripts/        # Helper scripts
└── results/        # Benchmark results
```

## Adding a New Implementation

1. Create folder: `{language}-{library}/`
2. Implement benchmark following existing examples
3. Output format: `{pattern_name} {time_ms} ms {count} matches`
4. Submit PR

## License

MIT
