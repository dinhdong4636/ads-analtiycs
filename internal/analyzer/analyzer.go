package analyzer

import (
	"database/sql"
	"fmt"

	"github.com/dinhdong4636/ads-analytics/internal/models"
	_ "github.com/marcboeker/go-duckdb"
)

type Analyzer struct {
	db *sql.DB
}

func New() (*Analyzer, error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, fmt.Errorf("failed to open DuckDB: %w", err)
	}

	return &Analyzer{db: db}, nil
}

func (a *Analyzer) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

func (a *Analyzer) AnalyzeCSV(csvPath string) (*models.AnalysisResult, error) {
	topCTR, err := a.getTopCTR(csvPath, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top CTR: %w", err)
	}

	topCPA, err := a.getTopCPA(csvPath, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top CPA: %w", err)
	}

	return &models.AnalysisResult{
		TopCTR: topCTR,
		TopCPA: topCPA,
	}, nil
}

func (a *Analyzer) getTopCTR(csvPath string, limit int) ([]models.CampaignMetrics, error) {
	query := fmt.Sprintf(`
		SELECT 
			campaign_id,
			SUM(impressions) as total_impressions,
			SUM(clicks) as total_clicks,
			SUM(spend) as total_spend,
			SUM(conversions) as total_conversions,
			CAST(SUM(clicks) AS DOUBLE) / NULLIF(SUM(impressions), 0) as ctr,
			CASE 
				WHEN SUM(conversions) > 0 THEN SUM(spend) / SUM(conversions)
				ELSE NULL
			END as cpa
		FROM read_csv_auto('%s')
		GROUP BY campaign_id
		HAVING SUM(impressions) > 0
		ORDER BY ctr DESC
		LIMIT %d
	`, csvPath, limit)

	return a.executeQuery(query)
}

func (a *Analyzer) getTopCPA(csvPath string, limit int) ([]models.CampaignMetrics, error) {
	query := fmt.Sprintf(`
		SELECT 
			campaign_id,
			SUM(impressions) as total_impressions,
			SUM(clicks) as total_clicks,
			SUM(spend) as total_spend,
			SUM(conversions) as total_conversions,
			CAST(SUM(clicks) AS DOUBLE) / NULLIF(SUM(impressions), 0) as ctr,
			SUM(spend) / SUM(conversions) as cpa
		FROM read_csv_auto('%s')
		GROUP BY campaign_id
		HAVING SUM(conversions) > 0
		ORDER BY cpa ASC
		LIMIT %d
	`, csvPath, limit)

	return a.executeQuery(query)
}

func (a *Analyzer) executeQuery(query string) ([]models.CampaignMetrics, error) {
	rows, err := a.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var results []models.CampaignMetrics
	for rows.Next() {
		var m models.CampaignMetrics
		err := rows.Scan(
			&m.CampaignID,
			&m.TotalImpressions,
			&m.TotalClicks,
			&m.TotalSpend,
			&m.TotalConversions,
			&m.CTR,
			&m.CPA,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}
