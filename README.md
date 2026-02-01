# Ads Analytics Service

A high-performance CLI tool for analyzing advertising campaign data from CSV/ZIP files, built with Go and DuckDB.

## Features
- **High Performance**: Powered by [DuckDB](https://duckdb.org/) for processing large datasets efficiently.
- **Support for ZIP**: Automatically extracts and processes CSV files directly from ZIP archives.
- **Metrics Calculation**:
  - Calculates CTR (Click-Through Rate)
  - Calculates CPA (Cost Per Acquisition)
- **Data Aggregation**: Aggregates multi-day data for unique campaigns.
- **Export**: Generates reports (Top 10 CTR & Top 10 CPA) in CSV format.
- **Dockerized**: Fully containerized environment for consistent development and execution.

## Tech Stack
- **Language**: Go 1.23+
- **Database Engine**: DuckDB (In-memory OLAP)
- **Containerization**: Docker & Docker Compose
- **Key Libraries**:
  - `github.com/marcboeker/go-duckdb`: Go driver for DuckDB.

---

## Setup Instructions

### Prerequisites
- [Docker](https://www.docker.com/) & Docker Compose installed.

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/dinhdong4636/ads-analytics.git
   cd ads-analytics
   ```

2. Start the development environment:
   ```bash
   cd docker
   docker compose up -d --build
   ```
   *This will build the container and start `ads-analytics-dev`. Air is pre-installed and waits for code changes to automatically recompile the binary.*

---

## Usage

### Running the Analyzer
Since the project is set up with hot-reloading, the binary is automatically compiled to `bin/main`.

To run the analyzer, you can simply execute the binary inside the container:

```bash
docker compose exec analytics-dev ./bin/main -csv <path_to_your_file>
```

**Example:**
```bash
# View help and available options
docker compose exec analytics-dev ./bin/main --help

# Analyze the default file (ad_data.csv.zip) if present
docker compose exec analytics-dev ./bin/main

# Analyze a specific local CSV file
docker compose exec analytics-dev ./bin/main -csv data/large_dataset.csv

# Analyze a ZIP archive (automatically extracts CSV)
docker compose exec analytics-dev ./bin/main -csv data/large_dataset.csv.zip
```

### Development (Hot-Reload)
The container uses [Air](https://github.com/cosmtrek/air) for live reloading. 
- Try editing any `.go` file in your IDE.
- Watch `docker logs -f ads-analytics-dev` to see Air automatically detect changes and rebuild `./bin/main` in seconds.


### Output
The tool will:
1. Print the analysis summary to the console.
2. Export detailed reports to the `export/` directory:
   - `export/top_ctr_<timestamp>.csv`
   - `export/top_cpa_<timestamp>.csv`

---

## Testing & Benchmarks

The project includes built-in Unit Tests and Benchmarks to ensure correctness and measure performance.

### Run Unit Tests
Verifies the logic for CTR/CPA calculations and data sorting.

```bash
docker exec ads-analytics-dev go test -v ./internal/analyzer
```

**Real Output:**
```text
=== RUN   TestAnalyzeCSV
--- PASS: TestAnalyzeCSV (0.01s)
PASS
ok      github.com/dinhdong4636/ads-analytics/internal/analyzer 0.016s
```

### Run Benchmarks
Measures performance across different dataset sizes (10, 100, 1k, 10k rows).

```bash
docker exec ads-analytics-dev go test -bench=. -benchmem -v ./internal/analyzer
```

**Real Output:**
```text
goos: linux
goarch: amd64
pkg: github.com/dinhdong4636/ads-analytics/internal/analyzer
cpu: Intel(R) Core(TM) i5-10400 CPU @ 2.90GHz
BenchmarkAnalyzeCSV
BenchmarkAnalyzeCSV/Rows_10000
BenchmarkAnalyzeCSV/Rows_10000-12                     28          40421668 ns/op           17003 B/op        864 allocs/op
BenchmarkAnalyzeCSV/Rows_1000
BenchmarkAnalyzeCSV/Rows_1000-12                      27          43373889 ns/op           16581 B/op        863 allocs/op
BenchmarkAnalyzeCSV/Rows_100
BenchmarkAnalyzeCSV/Rows_100-12                      178           6606430 ns/op           16402 B/op        863 allocs/op
BenchmarkAnalyzeCSV/Rows_10
BenchmarkAnalyzeCSV/Rows_10-12                       420           2789723 ns/op           16224 B/op        843 allocs/op
PASS
ok      github.com/dinhdong4636/ads-analytics/internal/analyzer 7.276s
```

**Understanding Benchmark Output:**
- `ns/op`: Nanoseconds per operation (execution time).
- `B/op`: Bytes allocated per operation.
- `allocs/op`: Number of memory allocations per operation.

---

## Performance Report

### Processing Performance
*Measured on Intel Core i5-10400 CPU @ 2.90GHz inside Docker*

| Dataset Size | Execution Time | Memory Alloc | Notes |
| :--- | :--- | :--- | :--- |
| **10 Rows** | ~2.79 ms | ~16 KB | Instant |
| **100 Rows** | ~6.60 ms | ~16 KB | Instant |
| **1,000 Rows** | ~43.37 ms | ~16 KB | Fast overhead dominated |
| **10,000 Rows** | ~40.42 ms | ~17 KB | Fast efficient scaling |
| **1GB File** (Est.) | ~15-20s | Low RAM | DuckDB uses streaming engine |

### Memory Usage
- **Peak Memory**: The application is highly memory-efficient. DuckDB processes data in blocks (vectorized execution) rather than loading the entire file into RAM.
- **Allocation overhead**: ~16KB per query execution (mostly static overhead from the CGO/DuckDB driver interface).

### Optimization Analysis
Benchmarks show that processing time scales very efficiently. The initial overhead (~40ms) dominates small datasets, but for larger files (1GB+), the throughput is limited primarily by disk I/O, not CPU.

---
