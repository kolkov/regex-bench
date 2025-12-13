use regex::Regex;
use std::env;
use std::fs;
use std::time::Instant;

struct Pattern {
    name: &'static str,
    pattern: &'static str,
}

const PATTERNS: &[Pattern] = &[
    Pattern { name: "literal_alt", pattern: r"error|warning|fatal|critical" },
    Pattern { name: "anchored", pattern: r"^HTTP/[12]\.[01]" },
    Pattern { name: "inner_literal", pattern: r".*@example\.com" },
    Pattern { name: "suffix", pattern: r".*\.(txt|log|md)" },
    Pattern { name: "char_class", pattern: r"[\w]+" },
    Pattern { name: "email", pattern: r"[\w.+-]+@[\w.-]+\.[\w.-]+" },
    Pattern { name: "uri", pattern: r"[\w]+://[^/\s?#]+[^\s?#]+(?:\?[^\s#]*)?(?:#[^\s]*)?" },
    Pattern { name: "ip", pattern: r"(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])" },
];

fn measure(data: &str, p: &Pattern) {
    let start = Instant::now();

    let re = Regex::new(p.pattern).expect("Invalid regex");
    let count = re.find_iter(data).count();

    let elapsed = start.elapsed();
    let ms = elapsed.as_secs_f64() * 1000.0;

    println!("{:<15} {:>10.2} ms  {:>6} matches", p.name, ms, count);
}

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() != 2 {
        eprintln!("Usage: benchmark <input-file>");
        std::process::exit(1);
    }

    let data = fs::read_to_string(&args[1]).expect("Failed to read file");
    let size_mb = data.len() as f64 / 1024.0 / 1024.0;

    println!("Rust regex (input: {:.2} MB)", size_mb);
    println!("─────────────────────────────────────────");

    for p in PATTERNS {
        measure(&data, p);
    }
}
