Following up on the coregex discussion with @mvdan - here are updated benchmarks (v0.8.22) with fair comparison (all engines on Linux, same input):

| Pattern | Go stdlib | Go coregex | Rust regex | coregex vs stdlib |
|---------|-----------|------------|------------|-------------------|
| literal_alt | 421 ms | 34 ms | **6 ms** | **12x faster** |
| anchored | 0.15 ms | **0.04 ms** | 0.31 ms | **4x faster** |
| inner_literal | 215 ms | 2.3 ms | **0.7 ms** | **94x faster** |
| suffix | 182 ms | 2.1 ms | **1.3 ms** | **87x faster** |
| char_class | 580 ms | **29 ms** | 65 ms | **20x faster** |
| email | 221 ms | 2.2 ms | **1.6 ms** | **99x faster** |

**Key findings:**
- coregex is **12-99x faster than stdlib** on all tested patterns
- Beats Rust on `char_class` (2.2x faster) and `anchored` (7.8x faster)
- Rust wins on `literal_alt` (Aho-Corasick) and `inner_literal`

For those who can't wait for stdlib improvements, coregex provides a drop-in alternative today.

Benchmark repo with CI for reproducible results: https://github.com/kolkov/regex-bench
