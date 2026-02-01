# Large-Scale Advertising Data Analysis with Go & DuckDB

## Problem Description
We need to process and analyze a large CSV dataset (~1GB) containing advertising performance records.  
The data is **read-only** and used purely for **analytical purposes**.

Each row represents daily metrics for an advertising campaign.

### CSV Schema
- `campaign_id` (string)
- `date` (YYYY-MM-DD)
- `impressions` (integer)
- `clicks` (integer)
- `spend` (float, USD)
- `conversions` (integer)

---

## Objectives
1. Efficiently process a large CSV file without loading it fully into memory  
2. Aggregate data by `campaign_id` to compute:
   - total_impressions  
   - total_clicks  
   - total_spend  
   - total_conversions  
3. Compute derived metrics:
   - **CTR** = total_clicks / total_impressions  
   - **CPA** = total_spend / total_conversions  
     - If `conversions = 0`, CPA should be ignored or set to NULL
4. Generate two result sets:
   - Top 10 campaigns with the **highest CTR**
   - Top 10 campaigns with the **lowest CPA**
5. Optimize for performance, memory efficiency, and clarity

---

## Chosen Solution

### Technology Stack
- **Go** as the implementation language
- **DuckDB** as the analytical database engine

### Architecture
- Query the CSV file directly using DuckDB (no pre-ingestion step)
- Perform aggregation, metric calculation, and ranking using SQL
- Execute the analysis in a single analytical pipeline
- Return only the final Top-N results to the application layer

---

## Why DuckDB + Go

### DuckDB Advantages
- Designed specifically for analytical (OLAP) workloads
- Can query large CSV files directly using `read_csv_auto`
- Vectorized and columnar execution for fast aggregation
- Low memory usage with automatic batching and disk spilling
- SQL-based logic is concise, readable, and easy to maintain

### Go Advantages
- High performance and fast execution
- Simple deployment as a single binary
- Strong suitability for batch processing and data pipelines
- Clean integration with DuckDB as an embedded database

---

## Advantages Over Other Approaches

- **Compared to Pandas**  
  - No manual chunking logic required  
  - Lower memory usage  
  - Better performance for large aggregations  

- **Compared to Polars**  
  - SQL-based aggregation is clearer for analytical queries  
  - Easier to express Top-N and metric logic  

- **Compared to Spark**  
  - No cluster or infrastructure overhead  
  - Ideal for single-machine workloads  

- **Compared to Traditional Databases**  
  - No data ingestion or schema setup required  
  - No server to manage  

---

## Summary
This solution uses **Go + DuckDB** to efficiently process large, read-only analytical datasets.  
It minimizes memory usage, maximizes performance, and keeps the logic simple and maintainable by leveraging DuckDB’s SQL-based analytical engine directly on CSV data.
