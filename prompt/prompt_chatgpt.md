# Prompt with ChatGPT

> Phạm vi: Chatgpt Free + Research.
> Mục đích: Collect thông tin + research thêm thông tin, dựa trên thông tin đó quyết định hướng đi của prj. Cuối cùng đẩy thông tin ra làm prompt / training / reference  để dùng cho model khác nếu cần.

---

## 1. Người dùng

I have a large CSV dataset (~1GB) containing advertising performance records.  
The goal is to handle large datasets efficiently, optimize performance/memory usage, and design a robust data-processing workflow.

example data Schema:

CSV Schema  
Column | Type | Description  
campaign_id | string | Campaign ID  
date | string | Date in YYYY-MM-DD format  
impressions | integer | Number of impressions  
clicks | integer | Number of clicks  
spend | float | Advertising cost (USD)  
conversions | integer | Number of conversions  

Example:

campaign_id | date | impressions | clicks | spend | conversions  
CMP001 | 2025-01-01 | 12000 | 300 | 45.50 | 12  
CMP002 | 2025-01-01 | 8000 | 120 | 28.00 | 4  
CMP001 | 2025-01-02 | 14000 | 340 | 48.20 | 15  
CMP003 | 2025-01-01 | 5000 | 60 | 15.00 | 3  
CMP002 | 2025-01-02 | 8500 | 150 | 31.00 | 5  

==> I need a solution to process and analyze data as quickly as possible to generate two lists:  
- Top 10 campaigns with the highest CTR  
- Top 10 campaigns with the lowest CPA  

Aggregate data by campaign_id.

For each campaign_id, compute:
- total_impressions  
- total_clicks  
- total_spend  
- total_conversions  
CTR = total_clicks / total_impressions  
CPA = total_spend / total_conversions  
If conversions = 0, ignore or return null for CPA.

---

## 2. Người dùng

First, I need a draft solution, for example, the architecture to be used, the technologies to be used, and the reasons for using them.

---

## 3. Người dùng

Okay, I understand that data that needs analysis doesn't require much manipulation to change it.  
Therefore, I think there's a type of database, or perhaps a data frame, called DuckDB.  
Given the file size, I think using it is reasonable.  
What do you think?

---

## 4. Người dùng

Okay please tell me the advantages and disadvantages of using DuckDB.  
What are its advantages over pandas and other types?

---

## 5. Người dùng

Okay, I've decided to use Go and DuckDB as the key to solve this problem.

Now, let's summarize this exchange into a prompt, as simple as possible but still retaining the main points.  

The prompt should include information such as:
- the requirements
- the objectives
- the solution
- the advantages of the solution over other approaches  

I need this prompt file to use in other models and so that everyone who reads it can understand the information.

---

## 6. Người dùng

I need a file md
